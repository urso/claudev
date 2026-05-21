package scan

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ErrNotGitRepo is returned by Churn when root is not inside a git work tree.
var ErrNotGitRepo = errors.New("not a git repository")

// ChurnSignals is the raw output of Churn — git log --numstat aggregation
// over a time window. Goal: code-ownership signal for onboarding
// ("who can I ask about X?").
type ChurnSignals struct {
	Since        string         `json:"since" jsonschema:"description=the --since value as passed (e.g. '90d', '6m', '1y')"`
	WindowStart  string         `json:"window_start" jsonschema:"description=RFC3339 timestamp of the start of the window"`
	Note         string         `json:"note,omitempty" jsonschema:"description=optional human-readable note (e.g. 'no commits in window' for shallow clones)"`
	TopFiles     []FileChurn    `json:"top_files" jsonschema:"description=top N files by lines changed (insertions+deletions) in window; sorted desc"`
	Contributors []Contributor  `json:"contributors" jsonschema:"description=global per-author aggregate sorted by lines desc; high lines_recent vs lines = currently active maintainer"`
	DirOwnership []DirOwnership `json:"dir_ownership" jsonschema:"description=per-top-level-dir ownership; answers 'who owns code under <dir>?'"`
	LargeCommits []CommitInfo   `json:"large_commits" jsonschema:"description=commits whose total touched lines >= --large-min-lines (default 500); capped at top 20"`
}

type FileChurn struct {
	Path       string `json:"path"`
	Commits    int    `json:"commits"`
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
}

// Contributor is a per-author churn aggregate. Lines = insertions+deletions.
// LinesRecent covers the most recent third of the window — high ratio vs
// Lines indicates a currently-active maintainer.
type Contributor struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Commits     int     `json:"commits"`
	Lines       int     `json:"lines" jsonschema:"description=insertions + deletions in window"`
	LinesRecent int     `json:"lines_recent" jsonschema:"description=lines changed in most-recent third of window; 0 = inactive/departed"`
	Share       float64 `json:"share,omitempty" jsonschema:"description=fraction of total churn lines (global) or dir lines (within dir_ownership)"`
	LastCommit  string  `json:"last_commit" jsonschema:"description=RFC3339 timestamp of author's most recent commit in window"`
}

// DirOwnership lists the top contributors for a top-level directory.
type DirOwnership struct {
	Dir     string        `json:"dir"`
	Lines   int           `json:"lines" jsonschema:"description=total ins+del across all authors for this dir"`
	Authors []Contributor `json:"authors" jsonschema:"description=top N authors by lines; share = author_lines / dir.lines"`
}

type CommitInfo struct {
	SHA     string `json:"sha"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Subject string `json:"subject"`
	Lines   int    `json:"lines" jsonschema:"description=total touched lines (ins+del) across all files"`
}

// ChurnOptions tunes Churn.
type ChurnOptions struct {
	Since               string // e.g. "90d", "30d", "6m"; default "90d"
	TopFilesN           int    // default 25
	TopContributorsN    int    // default 25
	TopAuthorsPerDirN   int    // default 3
	LargeCommitMinLines int    // default 500
}

func (o *ChurnOptions) applyDefaults() {
	if o.Since == "" {
		o.Since = "90d"
	}
	if o.TopFilesN <= 0 {
		o.TopFilesN = 25
	}
	if o.TopContributorsN <= 0 {
		o.TopContributorsN = 25
	}
	if o.TopAuthorsPerDirN <= 0 {
		o.TopAuthorsPerDirN = 3
	}
	if o.LargeCommitMinLines <= 0 {
		o.LargeCommitMinLines = 500
	}
}

// Churn parses git log --numstat over the window and aggregates signals.
//
// Pipeline: defaults → git probe → run log → parse aggregate → finalize.
// The parse and finalize stages are pure functions of their inputs and
// covered by their own tests.
func Churn(root string, opts ChurnOptions) (ChurnSignals, error) {
	opts.applyDefaults()

	if err := exec.Command("git", "-C", root, "rev-parse", "--is-inside-work-tree").Run(); err != nil {
		return ChurnSignals{}, ErrNotGitRepo
	}

	dur, err := parseDuration(opts.Since)
	if err != nil {
		return ChurnSignals{}, err
	}
	now := time.Now()
	windowStart := now.Add(-dur)
	recentCutoff := now.Add(-dur / 3)

	data, err := runGitLog(root, opts.Since)
	if err != nil {
		return ChurnSignals{}, err
	}

	agg, err := parseChurnLog(data, recentCutoff, opts.LargeCommitMinLines)
	if err != nil {
		return ChurnSignals{}, err
	}

	sig := finalize(agg, opts)
	sig.Since = opts.Since
	sig.WindowStart = windowStart.UTC().Format(time.RFC3339)
	if len(sig.Contributors) == 0 && len(sig.TopFiles) == 0 {
		sig.Note = "no commits in window (shallow clone or empty history?)"
	}
	return sig, nil
}

// runGitLog shells out to `git log --numstat` and returns its stdout.
func runGitLog(root, since string) ([]byte, error) {
	cmd := exec.Command("git", "-C", root, "log",
		"--since="+gitSince(since),
		"--numstat",
		"--pretty=format:%x1ecommit%x09%H%x09%an%x09%ae%x09%aI%x09%s",
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}
	return out, nil
}

// churnAggregate is the raw parse output, pre-finalization.
type churnAggregate struct {
	files        map[string]*fileAgg
	contribs     map[string]*Contributor
	dirAuthors   map[dirAuthorKey]*dirAuthorAgg
	dirTotals    map[string]int
	largeCommits []CommitInfo
}

type fileAgg struct {
	commits map[string]struct{}
	ins     int
	del     int
}

type dirAuthorKey struct{ dir, email string }

type dirAuthorAgg struct {
	name        string
	email       string
	lines       int
	linesRecent int
	commits     int
	lastCommit  string
}

// parseChurnLog parses git-log --numstat output into an aggregate. Pure
// function: no I/O. Commit lines are prefixed with \x1e to disambiguate
// from numstat rows.
func parseChurnLog(data []byte, recentCutoff time.Time, largeMin int) (churnAggregate, error) {
	agg := churnAggregate{
		files:      map[string]*fileAgg{},
		contribs:   map[string]*Contributor{},
		dirAuthors: map[dirAuthorKey]*dirAuthorAgg{},
		dirTotals:  map[string]int{},
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)

	var (
		curSHA, curName, curEmail, curDate, curSubject string
		curIsRecent                                    bool
		curLines                                       int
	)

	flushLarge := func() {
		if curSHA == "" || curLines < largeMin {
			return
		}
		agg.largeCommits = append(agg.largeCommits, CommitInfo{
			SHA: curSHA, Author: curName, Date: curDate,
			Subject: curSubject, Lines: curLines,
		})
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "\x1ecommit\t") {
			flushLarge()
			parts := strings.SplitN(line[len("\x1ecommit\t"):], "\t", 5)
			if len(parts) < 5 {
				continue
			}
			curSHA, curName, curEmail, curDate, curSubject = parts[0], parts[1], parts[2], parts[3], parts[4]
			curLines = 0

			if t, terr := time.Parse(time.RFC3339, curDate); terr == nil {
				curIsRecent = !t.Before(recentCutoff)
			} else {
				curIsRecent = false
			}

			c, ok := agg.contribs[curEmail]
			if !ok {
				c = &Contributor{Name: curName, Email: curEmail}
				agg.contribs[curEmail] = c
			}
			c.Commits++
			if c.LastCommit == "" || curDate > c.LastCommit {
				c.LastCommit = curDate
			}
			continue
		}
		if curSHA == "" {
			continue
		}
		ins, del, path, ok := parseNumstatLine(line)
		if !ok {
			continue
		}
		lines := ins + del
		curLines += lines
		recordFile(agg.files, path, curSHA, ins, del)
		if c := agg.contribs[curEmail]; c != nil {
			c.Lines += lines
			if curIsRecent {
				c.LinesRecent += lines
			}
		}
		recordDirAuthor(agg, path, curName, curEmail, curDate, lines, curIsRecent)
	}
	flushLarge()
	if err := scanner.Err(); err != nil {
		return churnAggregate{}, err
	}
	return agg, nil
}

// parseNumstatLine parses an "ins<TAB>del<TAB>path" line. Binary files
// have "-" for counts which Atoi treats as 0; that is the desired result.
func parseNumstatLine(line string) (ins, del int, path string, ok bool) {
	fields := strings.SplitN(line, "\t", 3)
	if len(fields) != 3 || fields[2] == "" {
		return 0, 0, "", false
	}
	ins, _ = strconv.Atoi(fields[0])
	del, _ = strconv.Atoi(fields[1])
	return ins, del, fields[2], true
}

func recordFile(files map[string]*fileAgg, path, sha string, ins, del int) {
	fa, ok := files[path]
	if !ok {
		fa = &fileAgg{commits: map[string]struct{}{}}
		files[path] = fa
	}
	fa.commits[sha] = struct{}{}
	fa.ins += ins
	fa.del += del
}

func recordDirAuthor(agg churnAggregate, path, name, email, date string, lines int, recent bool) {
	dir := topDirOf(path)
	agg.dirTotals[dir] += lines
	key := dirAuthorKey{dir: dir, email: email}
	da, ok := agg.dirAuthors[key]
	if !ok {
		da = &dirAuthorAgg{name: name, email: email}
		agg.dirAuthors[key] = da
	}
	da.lines += lines
	if recent {
		da.linesRecent += lines
	}
	da.commits++
	if da.lastCommit == "" || date > da.lastCommit {
		da.lastCommit = date
	}
}

// finalize sorts, caps, and computes share fields from a parsed aggregate.
func finalize(agg churnAggregate, opts ChurnOptions) ChurnSignals {
	return ChurnSignals{
		TopFiles:     topFiles(agg.files, opts.TopFilesN),
		Contributors: rankContributors(agg.contribs, opts.TopContributorsN),
		DirOwnership: dirOwnership(agg.dirAuthors, agg.dirTotals, opts.TopAuthorsPerDirN),
		LargeCommits: capLargeCommits(agg.largeCommits),
	}
}

func topFiles(files map[string]*fileAgg, n int) []FileChurn {
	out := make([]FileChurn, 0, len(files))
	for p, fa := range files {
		out = append(out, FileChurn{
			Path: p, Commits: len(fa.commits),
			Insertions: fa.ins, Deletions: fa.del,
		})
	}
	slices.SortFunc(out, func(a, b FileChurn) int {
		return cmp.Or(
			cmp.Compare(b.Insertions+b.Deletions, a.Insertions+a.Deletions),
			cmp.Compare(a.Path, b.Path),
		)
	})
	if len(out) > n {
		out = out[:n]
	}
	return out
}

func rankContributors(contribs map[string]*Contributor, n int) []Contributor {
	totalLines := 0
	out := make([]Contributor, 0, len(contribs))
	for _, c := range contribs {
		totalLines += c.Lines
		out = append(out, *c)
	}
	for i := range out {
		if totalLines > 0 {
			out[i].Share = float64(out[i].Lines) / float64(totalLines)
		}
	}
	slices.SortFunc(out, func(a, b Contributor) int {
		return cmp.Or(cmp.Compare(b.Lines, a.Lines), cmp.Compare(a.Email, b.Email))
	})
	if len(out) > n {
		out = out[:n]
	}
	return out
}

func dirOwnership(dirAuthors map[dirAuthorKey]*dirAuthorAgg, dirTotals map[string]int, perDirN int) []DirOwnership {
	groups := map[string][]Contributor{}
	for k, da := range dirAuthors {
		groups[k.dir] = append(groups[k.dir], Contributor{
			Name: da.name, Email: da.email,
			Commits: da.commits, Lines: da.lines,
			LinesRecent: da.linesRecent, LastCommit: da.lastCommit,
		})
	}
	out := make([]DirOwnership, 0, len(groups))
	for dir, authors := range groups {
		slices.SortFunc(authors, func(a, b Contributor) int {
			return cmp.Or(cmp.Compare(b.Lines, a.Lines), cmp.Compare(a.Email, b.Email))
		})
		total := dirTotals[dir]
		for i := range authors {
			if total > 0 {
				authors[i].Share = float64(authors[i].Lines) / float64(total)
			}
		}
		if len(authors) > perDirN {
			authors = authors[:perDirN]
		}
		out = append(out, DirOwnership{Dir: dir, Lines: total, Authors: authors})
	}
	slices.SortFunc(out, func(a, b DirOwnership) int {
		return cmp.Or(cmp.Compare(b.Lines, a.Lines), cmp.Compare(a.Dir, b.Dir))
	})
	return out
}

func capLargeCommits(large []CommitInfo) []CommitInfo {
	slices.SortFunc(large, func(a, b CommitInfo) int {
		return cmp.Or(cmp.Compare(b.Lines, a.Lines), cmp.Compare(b.Date, a.Date))
	})
	const cap = 20
	if len(large) > cap {
		large = large[:cap]
	}
	return large
}

func topDirOf(path string) string {
	head, _, ok := strings.Cut(path, "/")
	if !ok {
		return "(root)"
	}
	return head
}

// gitSince translates our shorthand into a string git's --since accepts.
// Git understands "N days ago" / "N months ago" / "N years ago" but not
// our compact "Nd/Nm/Ny" form.
func gitSince(s string) string {
	if s == "" {
		return s
	}
	last := s[len(s)-1]
	if _, err := strconv.Atoi(s[:len(s)-1]); err != nil {
		return s
	}
	switch last {
	case 'd':
		return s[:len(s)-1] + " days ago"
	case 'm':
		return s[:len(s)-1] + " months ago"
	case 'y':
		return s[:len(s)-1] + " years ago"
	}
	return s
}

// parseDuration accepts Go's time.ParseDuration plus shorthand "Nd" (days)
// and "Nm" (30-day months) / "Ny" (365-day years).
func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("empty duration")
	}
	last := s[len(s)-1]
	switch last {
	case 'd', 'm', 'y':
		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", s, err)
		}
		switch last {
		case 'd':
			return time.Duration(n) * 24 * time.Hour, nil
		case 'm':
			return time.Duration(n) * 30 * 24 * time.Hour, nil
		case 'y':
			return time.Duration(n) * 365 * 24 * time.Hour, nil
		}
	}
	return time.ParseDuration(s)
}

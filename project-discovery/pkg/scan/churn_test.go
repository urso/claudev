package scan

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func git(t *testing.T, dir string, env []string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	root := t.TempDir()
	git(t, root, nil, "init", "-q", "-b", "main")
	git(t, root, nil, "config", "user.email", "init@example.com")
	git(t, root, nil, "config", "user.name", "Init")
	git(t, root, nil, "config", "commit.gpgsign", "false")
	return root
}

func writeAndCommit(t *testing.T, root, rel, content, name, email, date, msg string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, nil, "add", rel)
	env := []string{
		"GIT_AUTHOR_NAME=" + name,
		"GIT_AUTHOR_EMAIL=" + email,
		"GIT_COMMITTER_NAME=" + name,
		"GIT_COMMITTER_EMAIL=" + email,
		"GIT_AUTHOR_DATE=" + date,
		"GIT_COMMITTER_DATE=" + date,
	}
	cmd := exec.Command("git", "commit", "-q", "-m", msg)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("commit: %v\n%s", err, out)
	}
}

func TestChurn_NotGitRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	root := t.TempDir()
	_, err := Churn(root, ChurnOptions{})
	if !errors.Is(err, ErrNotGitRepo) {
		t.Fatalf("err = %v, want ErrNotGitRepo", err)
	}
}

func TestChurn_AggregatesFilesAndContributors(t *testing.T) {
	root := initRepo(t)
	writeAndCommit(t, root, "pkg/a/a.go", "package a\n// 1\n// 2\n",
		"Alice", "alice@example.com", "2026-04-20T10:00:00Z", "init a")
	writeAndCommit(t, root, "pkg/a/a.go", "package a\n// 1\n// 2\n// 3\n// 4\n",
		"Alice", "alice@example.com", "2026-04-25T10:00:00Z", "expand a")
	writeAndCommit(t, root, "pkg/b/b.go", "package b\n",
		"Bob", "bob@example.com", "2026-04-28T10:00:00Z", "init b")

	sig, err := Churn(root, ChurnOptions{Since: "365d"})
	if err != nil {
		t.Fatal(err)
	}
	if len(sig.TopFiles) == 0 {
		t.Fatal("no top files")
	}
	if sig.TopFiles[0].Path != "pkg/a/a.go" {
		t.Errorf("top file = %q, want pkg/a/a.go", sig.TopFiles[0].Path)
	}
	if sig.TopFiles[0].Commits != 2 {
		t.Errorf("a.go commits = %d, want 2", sig.TopFiles[0].Commits)
	}

	byEmail := map[string]Contributor{}
	for _, c := range sig.Contributors {
		byEmail[c.Email] = c
	}
	if byEmail["alice@example.com"].Lines <= byEmail["bob@example.com"].Lines {
		t.Errorf("alice should outweigh bob: %+v", byEmail)
	}
	if byEmail["alice@example.com"].Commits != 2 {
		t.Errorf("alice commits = %d, want 2", byEmail["alice@example.com"].Commits)
	}

	dirByName := map[string]DirOwnership{}
	for _, d := range sig.DirOwnership {
		dirByName[d.Dir] = d
	}
	pkg, ok := dirByName["pkg"]
	if !ok {
		t.Fatalf("missing pkg dir ownership: %+v", sig.DirOwnership)
	}
	if pkg.Authors[0].Email != "alice@example.com" {
		t.Errorf("pkg top author = %q, want alice", pkg.Authors[0].Email)
	}
}

func TestParseDuration(t *testing.T) {
	cases := map[string]bool{
		"90d": true, "30d": true, "6m": true, "1y": true, "2h": true, "bogus": false,
	}
	for in, ok := range cases {
		_, err := parseDuration(in)
		if (err == nil) != ok {
			t.Errorf("parseDuration(%q) err=%v, ok=%v", in, err, ok)
		}
	}
}

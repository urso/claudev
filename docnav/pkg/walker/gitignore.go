package walker

import (
	"iter"

	gitignore "github.com/denormal/go-gitignore"
)

// FilterGitignore wraps a file iterator and filters out paths matched by
// .gitignore rules discovered from the repository root. If the gitignore
// rules cannot be loaded, the original iterator is returned unchanged.
func FilterGitignore(seq iter.Seq2[string, error]) iter.Seq2[string, error] {
	return FilterGitignoreAt(".", seq)
}

// FilterGitignoreAt wraps a file iterator and filters out paths matched by
// .gitignore rules discovered from the given repository base directory.
// If the gitignore rules cannot be loaded, the original iterator is returned unchanged.
func FilterGitignoreAt(base string, seq iter.Seq2[string, error]) iter.Seq2[string, error] {
	ignore, err := gitignore.NewRepository(base)
	if err != nil {
		return seq
	}

	return func(yield func(string, error) bool) {
		for path, err := range seq {
			if err != nil {
				if !yield("", err) {
					return
				}
				continue
			}
			if ignore.Ignore(path) {
				continue
			}
			if !yield(path, nil) {
				return
			}
		}
	}
}

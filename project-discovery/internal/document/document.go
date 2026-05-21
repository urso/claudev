package document

import (
	"bufio"
	"bytes"
	"errors"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/urso/claudev/project-discovery/internal/fsutil"
)

var ErrNoFrontmatter = errors.New("no frontmatter found")

// ErrNoPath is returned by Document.Write when called on a document that
// has no Path set (e.g. one constructed in memory rather than read from disk).
var ErrNoPath = errors.New("document has no path; use WriteFile")

var separator = []byte("---\n")

// Document represents a markdown file with YAML frontmatter.
//
// Path is the source location of the document. ParseDocumentFile sets it;
// in-memory construction leaves it empty until the caller assigns one.
// Write uses Path; WriteFile takes an explicit destination.
type Document[T any] struct {
	Frontmatter T
	Body        []byte
	Path        string
}

// ParseDocument parses a markdown document with YAML frontmatter into a Document[T].
func ParseDocument[T any](data []byte) (Document[T], error) {
	var doc Document[T]

	yamlData, body, err := extract(data)
	if err != nil {
		return doc, err
	}

	if err := yaml.Unmarshal(yamlData, &doc.Frontmatter); err != nil {
		return doc, err
	}
	doc.Body = body
	return doc, nil
}

// ParseDocumentFile reads a file and parses it as a Document[T].
// The returned document carries Path = path so it can be written back via Write.
func ParseDocumentFile[T any](path string) (Document[T], error) {
	data, err := os.ReadFile(path)
	if err != nil {
		var zero Document[T]
		return zero, err
	}
	doc, err := ParseDocument[T](data)
	if err != nil {
		return doc, err
	}
	doc.Path = path
	return doc, nil
}

// Bytes serializes the document back to markdown with YAML frontmatter.
// The body is preserved byte-for-byte.
func (d Document[T]) Bytes() ([]byte, error) {
	yamlData, err := yaml.Marshal(d.Frontmatter)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(separator)
	buf.Write(yamlData)
	buf.Write(separator)
	buf.Write(d.Body)
	return buf.Bytes(), nil
}

// WriteFile serializes the document and writes it atomically to path.
func (d Document[T]) WriteFile(path string) error {
	data, err := d.Bytes()
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(path, data)
}

// Write serializes the document and writes it atomically to its own Path.
// Returns ErrNoPath if the document has no Path set.
func (d Document[T]) Write() error {
	if d.Path == "" {
		return ErrNoPath
	}
	return d.WriteFile(d.Path)
}

// ParseFrontmatterFromFile opens a file and reads only the YAML frontmatter block,
// stopping at the closing "---" line. The file body is never read into memory.
func ParseFrontmatterFromFile[T any](path string) (T, error) {
	var zero T

	f, err := os.Open(path)
	if err != nil {
		return zero, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	if !scanner.Scan() || scanner.Text() != "---" {
		return zero, ErrNoFrontmatter
	}

	var buf bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			var fm T
			if err := yaml.Unmarshal(buf.Bytes(), &fm); err != nil {
				return zero, err
			}
			return fm, nil
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return zero, err
	}

	return zero, ErrNoFrontmatter
}

// extract splits a markdown document into the YAML frontmatter bytes and the body.
// The frontmatter is the content between the first two "---\n" lines.
func extract(data []byte) (yamlData []byte, body []byte, err error) {
	if !bytes.HasPrefix(data, separator) {
		return nil, nil, ErrNoFrontmatter
	}

	rest := data[len(separator):]
	idx := bytes.Index(rest, separator)
	if idx < 0 {
		return nil, nil, ErrNoFrontmatter
	}

	yamlData = rest[:idx]
	body = rest[idx+len(separator):]
	return yamlData, body, nil
}

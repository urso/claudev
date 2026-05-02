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

var separator = []byte("---\n")

// Document represents a markdown file with YAML frontmatter.
type Document[T any] struct {
	Frontmatter T
	Body        []byte
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
func ParseDocumentFile[T any](path string) (Document[T], error) {
	data, err := os.ReadFile(path)
	if err != nil {
		var zero Document[T]
		return zero, err
	}
	return ParseDocument[T](data)
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

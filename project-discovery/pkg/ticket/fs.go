package ticket

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urso/claudev/project-discovery/internal/document"
	"github.com/urso/claudev/project-discovery/internal/fsutil"
)

// DiscoveryDirName is the conventional directory name for discovery state in a repo.
const DiscoveryDirName = ".discovery"

// ErrDiscoveryDirNotFound is returned when no .discovery/ directory can be located.
var ErrDiscoveryDirNotFound = errors.New(".discovery/ not found")

// ErrTicketNotFound is returned when no ticket file matches a given id.
var ErrTicketNotFound = errors.New("ticket not found")

// Discovery represents a located .discovery/ directory in a repo.
type Discovery struct {
	root string
}

// FindDiscovery walks up from start looking for a .discovery/ directory,
// falling back to <git-root>/.discovery.
func FindDiscovery(start string) (Discovery, error) {
	root, err := fsutil.FindProjRoot(start, DiscoveryDirName)
	if err != nil {
		if errors.Is(err, fsutil.ErrNotFound) {
			return Discovery{}, ErrDiscoveryDirNotFound
		}
		return Discovery{}, err
	}
	return Discovery{root: root}, nil
}

// Path returns the absolute path of the .discovery/ directory.
func (d Discovery) Path() string { return d.root }

// Tickets returns a handle to the tickets/ subdirectory.
func (d Discovery) Tickets() TicketFolder {
	return TicketFolder{dir: filepath.Join(d.root, "tickets")}
}

// TicketFolder represents the tickets/ subdirectory of a Discovery.
type TicketFolder struct {
	dir string
}

// Path returns the absolute path of the tickets/ directory.
func (t TicketFolder) Path() string { return t.dir }

// resolve returns the absolute path of the ticket file whose frontmatter
// id equals id. Frontmatter is authoritative — filenames are only a hint.
func (t TicketFolder) resolve(id ID) (string, error) {
	hits, err := findByID(t, id)
	if err != nil {
		return "", err
	}
	switch len(hits) {
	case 0:
		return "", fmt.Errorf("%w: %s", ErrTicketNotFound, id)
	case 1:
		return hits[0], nil
	default:
		return "", fmt.Errorf("multiple tickets claim id %s: %v", id, hits)
	}
}

// Read loads the ticket with the given id, found via frontmatter lookup.
// The returned Document carries its source path; mutate and call doc.Write()
// to persist back. Returns ErrTicketNotFound if no file claims the id.
func (t TicketFolder) Read(id ID) (Document, error) {
	path, err := t.resolve(id)
	if err != nil {
		return Document{}, err
	}
	return document.ParseDocumentFile[Ticket](path)
}

// Write ensures the tickets directory exists and persists doc to its Path.
// Returns document.ErrNoPath if doc has no Path set.
func (t TicketFolder) Write(doc Document) error {
	if err := os.MkdirAll(t.dir, 0o755); err != nil {
		return err
	}
	return doc.Write()
}

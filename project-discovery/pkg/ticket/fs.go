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

// Read loads the ticket whose filename matches "<id>-*.md".
// Returns ErrTicketNotFound if no file matches, or an error if more than one matches.
func (t TicketFolder) Read(id string) (Document, error) {
	pattern := filepath.Join(t.dir, id+"-*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return Document{}, err
	}
	switch len(matches) {
	case 0:
		return Document{}, fmt.Errorf("%w: %s", ErrTicketNotFound, id)
	case 1:
		return document.ParseDocumentFile[Ticket](matches[0])
	default:
		return Document{}, fmt.Errorf("multiple tickets match id %q: %v", id, matches)
	}
}

// Write serializes doc and writes it atomically to path. The caller is
// responsible for choosing the path (e.g. "<id>-<slug>.md" inside Path()).
// Ensures the tickets directory exists.
func (t TicketFolder) Write(path string, doc Document) error {
	if err := os.MkdirAll(t.dir, 0755); err != nil {
		return err
	}
	return doc.WriteFile(path)
}

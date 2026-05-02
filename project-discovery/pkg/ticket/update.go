package ticket

import (
	"bytes"
	"fmt"
	"slices"
	"time"
)

// validStatuses lists the lifecycle values accepted by UpdateStatus.
var validStatuses = []string{StatusOpen, StatusActive, StatusDone, StatusDropped, StatusBlocked}

// UpdateStatus sets the status field of ticket id and bumps Updated.
// Rejects unknown status values.
func UpdateStatus(folder TicketFolder, id ID, status string) error {
	if !slices.Contains(validStatuses, status) {
		return fmt.Errorf("invalid status %q (want one of %v)", status, validStatuses)
	}
	doc, err := folder.Read(id)
	if err != nil {
		return err
	}
	doc.Frontmatter.Status = status
	doc.Frontmatter.Updated = today()
	return doc.Write()
}

// AppendLogs appends one log entry per message under the "## Log" section
// of ticket id. All entries share today's date and the file is written once.
// If the body has no "## Log" section it is created. Empty msgs is a no-op.
func AppendLogs(folder TicketFolder, id ID, msgs []string) error {
	if len(msgs) == 0 {
		return nil
	}
	doc, err := folder.Read(id)
	if err != nil {
		return err
	}
	date := today()
	doc.Body = appendLogEntries(doc.Body, date, msgs)
	doc.Frontmatter.Updated = date
	return doc.Write()
}

var logHeader = []byte("## Log")

// appendLogEntries inserts "- <date>: <msg>" lines into the Log section.
// New lines are appended after the last existing log content (or right after
// the header if the section is empty). The Log section is created at the end
// of the body when missing.
func appendLogEntries(body []byte, date string, msgs []string) []byte {
	var entries bytes.Buffer
	for _, m := range msgs {
		fmt.Fprintf(&entries, "- %s: %s\n", date, m)
	}

	idx := bytes.Index(body, logHeader)
	if idx < 0 {
		var out bytes.Buffer
		out.Write(body)
		if len(body) > 0 && !bytes.HasSuffix(body, []byte("\n")) {
			out.WriteByte('\n')
		}
		out.WriteString("\n## Log\n\n")
		out.Write(entries.Bytes())
		return out.Bytes()
	}

	// Find end of log section: next "## " heading at column 0, or EOF.
	tailStart := idx + len(logHeader)
	rest := body[tailStart:]
	end := len(body)
	for off := 0; ; {
		nl := bytes.IndexByte(rest[off:], '\n')
		if nl < 0 {
			break
		}
		lineStart := off + nl + 1
		if bytes.HasPrefix(rest[lineStart:], []byte("## ")) {
			end = tailStart + lineStart
			break
		}
		off = lineStart
	}

	// Trim trailing blank lines inside the section so new entries sit
	// flush against existing ones, then re-add a single trailing newline
	// before any following section.
	section := body[:end]
	suffix := body[end:]
	trimmed := bytes.TrimRight(section, "\n")

	var out bytes.Buffer
	out.Write(trimmed)
	out.WriteByte('\n')
	out.Write(entries.Bytes())
	if len(suffix) > 0 {
		out.WriteByte('\n')
		out.Write(suffix)
	}
	return out.Bytes()
}

func today() string { return time.Now().Format("2006-01-02") }

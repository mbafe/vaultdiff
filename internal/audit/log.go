package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/vaultdiff/internal/diff"
)

// Entry represents a single audit log record for a diff operation.
type Entry struct {
	Timestamp  time.Time        `json:"timestamp"`
	Path       string           `json:"path"`
	FromVersion int             `json:"from_version"`
	ToVersion   int             `json:"to_version"`
	Changes    []diff.Change    `json:"changes"`
	ChangeCount int             `json:"change_count"`
	User       string           `json:"user,omitempty"`
}

// Logger writes audit entries to an io.Writer.
type Logger struct {
	w      io.Writer
	format string
}

// NewLogger creates a Logger that writes to w in the given format ("json" or "text").
func NewLogger(w io.Writer, format string) *Logger {
	return &Logger{w: w, format: format}
}

// Record writes an audit entry for a completed diff operation.
func (l *Logger) Record(path string, fromVersion, toVersion int, changes []diff.Change, user string) error {
	entry := Entry{
		Timestamp:   time.Now().UTC(),
		Path:        path,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Changes:     changes,
		ChangeCount: len(changes),
		User:        user,
	}

	switch l.format {
	case "json":
		return l.writeJSON(entry)
	default:
		return l.writeText(entry)
	}
}

func (l *Logger) writeJSON(e Entry) error {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(l.w, string(data))
	return err
}

func (l *Logger) writeText(e Entry) error {
	_, err := fmt.Fprintf(l.w, "[%s] path=%s versions=%d->%d changes=%d user=%s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Path,
		e.FromVersion,
		e.ToVersion,
		e.ChangeCount,
		e.User,
	)
	return err
}

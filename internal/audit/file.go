package audit

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileLogger wraps Logger with a backing file, managing open/close lifecycle.
type FileLogger struct {
	*Logger
	f *os.File
}

// NewFileLogger opens (or creates/appends) the file at path and returns a FileLogger.
// format is passed through to Logger ("json" or "text").
func NewFileLogger(path, format string) (*FileLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file %q: %w", path, err)
	}

	return &FileLogger{
		Logger: NewLogger(f, format),
		f:      f,
	}, nil
}

// Close flushes and closes the underlying file.
func (fl *FileLogger) Close() error {
	if fl.f == nil {
		return nil
	}
	if err := fl.f.Sync(); err != nil {
		_ = fl.f.Close()
		return fmt.Errorf("audit: sync log file: %w", err)
	}
	return fl.f.Close()
}

package audit_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/audit"
)

func TestNewFileLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logs", "audit.log")

	fl, err := audit.NewFileLogger(path, "text")
	if err != nil {
		t.Fatalf("NewFileLogger error: %v", err)
	}
	defer fl.Close()

	if err := fl.Record("secret/test", 1, 2, nil, "tester"); err != nil {
		t.Fatalf("Record error: %v", err)
	}
	fl.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), "secret/test") {
		t.Errorf("log file missing expected content; got: %s", data)
	}
}

func TestNewFileLogger_AppendsBetweenSessions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	for i := 0; i < 3; i++ {
		fl, err := audit.NewFileLogger(path, "text")
		if err != nil {
			t.Fatalf("session %d: NewFileLogger error: %v", i, err)
		}
		_ = fl.Record("secret/app", i, i+1, nil, "")
		fl.Close()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 log lines, got %d:\n%s", len(lines), data)
	}
}

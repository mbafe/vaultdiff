package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/audit"
	"github.com/vaultdiff/internal/diff"
)

var sampleChanges = []diff.Change{
	{Key: "db_password", Type: diff.Modified},
	{Key: "api_key", Type: diff.Added},
}

func TestLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf, "text")

	if err := l.Record("secret/myapp", 1, 2, sampleChanges, "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{"secret/myapp", "1->2", "changes=2", "user=alice"} {
		if !strings.Contains(output, want) {
			t.Errorf("text output missing %q; got: %s", want, output)
		}
	}
}

func TestLogger_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf, "json")

	if err := l.Record("secret/myapp", 3, 4, sampleChanges, "bob"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if entry.Path != "secret/myapp" {
		t.Errorf("path = %q, want %q", entry.Path, "secret/myapp")
	}
	if entry.ChangeCount != 2 {
		t.Errorf("change_count = %d, want 2", entry.ChangeCount)
	}
	if entry.User != "bob" {
		t.Errorf("user = %q, want %q", entry.User, "bob")
	}
}

func TestLogger_EmptyChanges(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf, "json")

	if err := l.Record("secret/empty", 1, 1, nil, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.ChangeCount != 0 {
		t.Errorf("expected 0 changes, got %d", entry.ChangeCount)
	}
}

package diff

import (
	"strings"
	"testing"
)

func TestRender_NoChanges(t *testing.T) {
	out := Render(nil, RenderOptions{})
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected no-changes message, got: %q", out)
	}
}

func TestRender_AddedKey(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Key: "token", NewValue: "abc123"},
	}
	out := Render(changes, RenderOptions{MaskSecrets: false, Colorize: false})
	if !strings.Contains(out, "+ token = abc123") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_RemovedKey(t *testing.T) {
	changes := []Change{
		{Type: ChangeRemoved, Key: "old_key", OldValue: "val"},
	}
	out := Render(changes, RenderOptions{Colorize: false})
	if !strings.Contains(out, "- old_key = val") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_ModifiedKey(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Key: "db_pass", OldValue: "old", NewValue: "new"},
	}
	out := Render(changes, RenderOptions{Colorize: false})
	if !strings.Contains(out, "~ db_pass: old => new") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_MaskSecrets(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Key: "api_key", NewValue: "supersecret"},
	}
	out := Render(changes, RenderOptions{MaskSecrets: true, Colorize: false})
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected masked value, got full secret in: %q", out)
	}
	if !strings.Contains(out, "sup") {
		t.Errorf("expected visible prefix 'sup' in: %q", out)
	}
}

func TestRenderHeader(t *testing.T) {
	out := RenderHeader("secret/myapp", 1, 2, false)
	if !strings.Contains(out, "secret/myapp") {
		t.Errorf("expected path in header: %q", out)
	}
	if !strings.Contains(out, "v1") || !strings.Contains(out, "v2") {
		t.Errorf("expected version numbers in header: %q", out)
	}
}

func TestMaskValue_Short(t *testing.T) {
	if maskValue("") != "" {
		t.Error("expected empty string for empty input")
	}
	result := maskValue("ab")
	if result != "ab" {
		t.Errorf("short value should show fully, got: %q", result)
	}
}

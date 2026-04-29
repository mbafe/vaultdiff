package diff

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func makeVersion(ver int, data map[string]interface{}) *vault.SecretVersion {
	return &vault.SecretVersion{Version: ver, Data: data}
}

func TestCompare_NoChanges(t *testing.T) {
	old := makeVersion(1, map[string]interface{}{"key": "val"})
	new := makeVersion(2, map[string]interface{}{"key": "val"})

	r := Compare("secret", "myapp/config", old, new)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompare_AddedKey(t *testing.T) {
	old := makeVersion(1, map[string]interface{}{})
	new := makeVersion(2, map[string]interface{}{"token": "abc"})

	r := Compare("secret", "myapp/config", old, new)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(r.Changes) != 1 || r.Changes[0].Change != Added {
		t.Errorf("expected Added, got %v", r.Changes)
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	old := makeVersion(1, map[string]interface{}{"token": "abc"})
	new := makeVersion(2, map[string]interface{}{})

	r := Compare("secret", "myapp/config", old, new)
	if r.Changes[0].Change != Removed {
		t.Errorf("expected Removed, got %v", r.Changes[0].Change)
	}
}

func TestCompare_ModifiedKey(t *testing.T) {
	old := makeVersion(1, map[string]interface{}{"pass": "old"})
	new := makeVersion(2, map[string]interface{}{"pass": "new"})

	r := Compare("secret", "myapp/config", old, new)
	if r.Changes[0].Change != Modified {
		t.Errorf("expected Modified, got %v", r.Changes[0].Change)
	}
	if r.Changes[0].OldValue != "old" || r.Changes[0].NewValue != "new" {
		t.Errorf("unexpected old/new values: %v", r.Changes[0])
	}
}

func TestUnionKeys_Sorted(t *testing.T) {
	a := map[string]interface{}{"z": 1, "a": 2}
	b := map[string]interface{}{"m": 3}
	keys := unionKeys(a, b)
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

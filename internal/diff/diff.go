package diff

import (
	"fmt"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// ChangeType indicates the kind of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// KeyDiff represents a diff entry for a single secret key.
type KeyDiff struct {
	Key      string
	Change   ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two secret versions.
type Result struct {
	Mount      string
	Path       string
	OldVersion int
	NewVersion int
	Changes    []KeyDiff
}

// HasChanges returns true if any key was added, removed, or modified.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Change != Unchanged {
			return true
		}
	}
	return false
}

// Compare produces a diff Result between two SecretVersions.
func Compare(mount, path string, oldSV, newSV *vault.SecretVersion) *Result {
	result := &Result{
		Mount:      mount,
		Path:       path,
		OldVersion: oldSV.Version,
		NewVersion: newSV.Version,
	}

	allKeys := unionKeys(oldSV.Data, newSV.Data)
	for _, key := range allKeys {
		oldVal, oldOk := oldSV.Data[key]
		newVal, newOk := newSV.Data[key]

		var kd KeyDiff
		kd.Key = key

		switch {
		case !oldOk:
			kd.Change = Added
			kd.NewValue = fmt.Sprintf("%v", newVal)
		case !newOk:
			kd.Change = Removed
			kd.OldValue = fmt.Sprintf("%v", oldVal)
		case fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal):
			kd.Change = Modified
			kd.OldValue = fmt.Sprintf("%v", oldVal)
			kd.NewValue = fmt.Sprintf("%v", newVal)
		default:
			kd.Change = Unchanged
		}
		result.Changes = append(result.Changes, kd)
	}
	return result
}

// unionKeys returns a sorted slice of all keys from both maps.
func unionKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

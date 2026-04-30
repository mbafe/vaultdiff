package diff

// ChangeType represents the kind of difference detected between two secret versions.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// Change describes a single key-level difference between two secret versions.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Compare returns the list of Changes between two flat secret data maps.
// versionA is treated as the "before" state; versionB as "after".
func Compare(versionA, versionB map[string]string) []Change {
	var changes []Change

	for _, key := range unionKeys(versionA, versionB) {
		oldVal, inA := versionA[key]
		newVal, inB := versionB[key]

		switch {
		case inA && !inB:
			changes = append(changes, Change{
				Key:      key,
				Type:     ChangeRemoved,
				OldValue: oldVal,
			})
		case !inA && inB:
			changes = append(changes, Change{
				Key:      key,
				Type:     ChangeAdded,
				NewValue: newVal,
			})
		case inA && inB && oldVal != newVal:
			changes = append(changes, Change{
				Key:      key,
				Type:     ChangeModified,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	return changes
}

// unionKeys returns a sorted, deduplicated list of all keys present in either map.
func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
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

	// stable sort for deterministic output
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}

	return keys
}

package diff

import "strings"

// FilterOptions controls which changes are included in a filtered result.
type FilterOptions struct {
	// OnlyKeys, if non-empty, restricts output to changes whose key matches
	// one of the provided prefixes (case-insensitive).
	OnlyKeys []string

	// Types restricts output to specific change types. An empty slice means
	// all change types are included.
	Types []ChangeType
}

// Filter returns a new slice of Changes that satisfy the given FilterOptions.
// The original slice is not modified.
func Filter(changes []Change, opts FilterOptions) []Change {
	result := make([]Change, 0, len(changes))

	for _, c := range changes {
		if !matchesKeyFilter(c.Key, opts.OnlyKeys) {
			continue
		}
		if !matchesTypeFilter(c.Type, opts.Types) {
			continue
		}
		result = append(result, c)
	}

	return result
}

// matchesKeyFilter returns true when the key matches at least one prefix in
// prefixes, or when prefixes is empty (no filter applied).
func matchesKeyFilter(key string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true
	}
	lower := strings.ToLower(key)
	for _, p := range prefixes {
		if strings.HasPrefix(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// matchesTypeFilter returns true when the change type is in the allowed set,
// or when the set is empty (no filter applied).
func matchesTypeFilter(ct ChangeType, types []ChangeType) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if ct == t {
			return true
		}
	}
	return false
}

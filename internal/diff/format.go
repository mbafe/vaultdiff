package diff

import (
	"fmt"
	"io"
	"strings"
)

// Format controls how the diff output is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Render writes a human-readable diff to w.
func Render(w io.Writer, r *Result) {
	fmt.Fprintf(w, "--- %s/%s (v%d)\n", r.Mount, r.Path, r.OldVersion)
	fmt.Fprintf(w, "+++ %s/%s (v%d)\n", r.Mount, r.Path, r.NewVersion)

	if !r.HasChanges() {
		fmt.Fprintln(w, "(no changes)")
		return
	}

	for _, c := range r.Changes {
		switch c.Change {
		case Added:
			fmt.Fprintf(w, "+ %-30s = %s\n", c.Key, maskValue(c.NewValue))
		case Removed:
			fmt.Fprintf(w, "- %-30s = %s\n", c.Key, maskValue(c.OldValue))
		case Modified:
			fmt.Fprintf(w, "~ %-30s : %s -> %s\n", c.Key, maskValue(c.OldValue), maskValue(c.NewValue))
		}
	}
}

// maskValue replaces the content of a secret value with asterisks,
// preserving only the length hint.
func maskValue(v string) string {
	if len(v) == 0 {
		return `""`
	}
	return "[" + strings.Repeat("*", min(len(v), 8)) + "]"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

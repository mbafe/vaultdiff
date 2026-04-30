package diff

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	addedColor    = color.New(color.FgGreen)
	removedColor  = color.New(color.FgRed)
	modifiedColor = color.New(color.FgYellow)
	headerColor   = color.New(color.FgCyan, color.Bold)
)

// RenderOptions controls how the diff output is formatted.
type RenderOptions struct {
	MaskSecrets bool
	Colorize    bool
}

// Render formats a slice of Change records into a human-readable diff string.
func Render(changes []Change, opts RenderOptions) string {
	if len(changes) == 0 {
		return "No changes detected.\n"
	}

	var sb strings.Builder

	for _, c := range changes {
		switch c.Type {
		case ChangeAdded:
			line := fmt.Sprintf("+ %s = %s\n", c.Key, displayValue(c.NewValue, opts.MaskSecrets))
			if opts.Colorize {
				addedColor.Fprint(&sb, line)
			} else {
				sb.WriteString(line)
			}
		case ChangeRemoved:
			line := fmt.Sprintf("- %s = %s\n", c.Key, displayValue(c.OldValue, opts.MaskSecrets))
			if opts.Colorize {
				removedColor.Fprint(&sb, line)
			} else {
				sb.WriteString(line)
			}
		case ChangeModified:
			oldVal := displayValue(c.OldValue, opts.MaskSecrets)
			newVal := displayValue(c.NewValue, opts.MaskSecrets)
			line := fmt.Sprintf("~ %s: %s => %s\n", c.Key, oldVal, newVal)
			if opts.Colorize {
				modifiedColor.Fprint(&sb, line)
			} else {
				sb.WriteString(line)
			}
		}
	}

	return sb.String()
}

// RenderHeader prints a summary header for a diff operation.
func RenderHeader(path string, versionA, versionB int, colorize bool) string {
	line := fmt.Sprintf("=== %s (v%d → v%d) ===\n", path, versionA, versionB)
	if colorize {
		var sb strings.Builder
		headerColor.Fprint(&sb, line)
		return sb.String()
	}
	return line
}

func displayValue(v string, mask bool) string {
	if !mask {
		return v
	}
	return maskValue(v)
}

func maskValue(v string) string {
	if len(v) == 0 {
		return ""
	}
	visible := min(3, len(v))
	return v[:visible] + strings.Repeat("*", len(v)-visible)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

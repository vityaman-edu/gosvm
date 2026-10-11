// Package srcloc describes source text and source-located diagnostics.
package srcloc

import "fmt"

// Diag formats a source-located error, preserving fmt.Errorf's %w wrapping.
func Diag(path string, position Position, format string, args ...any) error {
	message := fmt.Errorf(format, args...) //nolint:err113 // Formatting helper: dynamic messages with %w wrapping.

	return fmt.Errorf(
		"%s:%d:%d: %w", path, position.Line, position.Character, message,
	)
}

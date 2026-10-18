package srcloc

import "fmt"

func Diag(path string, position Position, format string, args ...any) error {
	message := fmt.Errorf(format, args...) //nolint:err113 // Formatting helper: dynamic messages with %w wrapping.

	return fmt.Errorf(
		"%s:%d:%d: %w", path, position.Line, position.Character, message,
	)
}

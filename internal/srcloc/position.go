package srcloc

// Position identifies a location in source text using one-based coordinates.
type Position struct {
	Line      int
	Character int
}

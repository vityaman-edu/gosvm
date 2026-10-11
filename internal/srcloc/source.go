package srcloc

// FunctionName identifies a function to analyze within a source package.
type FunctionName string

// Source pairs a diagnostic file path with its Go source text.
type Source struct {
	Path string
	Text string
}

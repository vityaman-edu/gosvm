// Package gosmv provides the entry point for GoSVM program analysis.
package gosmv

import (
	"errors"
	"fmt"

	"github.com/vityaman-edu/gosvm/internal/gossa"
	"github.com/vityaman-edu/gosvm/internal/srcloc"
	"golang.org/x/tools/go/ssa"
)

// ErrFunctionNotFound indicates that the selected function is absent.
var ErrFunctionNotFound = errors.New("function not found")

// Analysis contains the SSA package used for subsequent program analysis.
type Analysis struct {
	Package *ssa.Package
}

// Analyze builds the source's SSA and validates the selected function.
func Analyze(
	source srcloc.Source, function srcloc.FunctionName,
) (Analysis, error) {
	pkg, err := gossa.BuildSSA(source)
	if err != nil {
		return Analysis{}, fmt.Errorf("build SSA for %q: %w", source.Path, err)
	}

	if pkg == nil {
		panic("BuildSSA returned a nil package without an error")
	}

	if pkg.Func(string(function)) == nil {
		return Analysis{}, fmt.Errorf(
			"function %q in %q: %w", function, source.Path, ErrFunctionNotFound,
		)
	}

	return Analysis{Package: pkg}, nil
}

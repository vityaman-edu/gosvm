package gossa

import (
	"fmt"
	"io"
	"slices"

	"golang.org/x/tools/go/ssa"
)

// PrintSSA writes the package and its function bodies in name order.
func PrintSSA(output io.Writer, pkg *ssa.Package) error {
	_, err := pkg.WriteTo(output)
	if err != nil {
		return fmt.Errorf("print SSA package: %w", err)
	}

	var names []string

	for name := range pkg.Members {
		names = append(names, name)
	}

	slices.Sort(names)

	for _, name := range names {
		function, ok := pkg.Members[name].(*ssa.Function)
		if !ok {
			continue
		}

		err := printFunction(output, function)
		if err != nil {
			return err
		}
	}

	return nil
}

func printFunction(output io.Writer, function *ssa.Function) error {
	_, err := function.WriteTo(output)
	if err != nil {
		return fmt.Errorf("print SSA function %q: %w", function.Name(), err)
	}

	for _, anonymous := range function.AnonFuncs {
		err := printFunction(output, anonymous)
		if err != nil {
			return err
		}
	}

	return nil
}

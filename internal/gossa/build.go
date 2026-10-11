// Package gossa builds and prints Go static single assignment (SSA) packages.
package gossa

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/vityaman-edu/gosvm/internal/srcloc"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// BuildSSA parses and type-checks a source file, then builds its SSA package.
func BuildSSA(source srcloc.Source) (*ssa.Package, error) {
	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(
		fileSet, source.Path, source.Text, parser.ParseComments,
	)
	if err != nil {
		return nil, fmt.Errorf("parse source: %w", err)
	}

	files := []*ast.File{file}
	pkg := types.NewPackage(source.Path, file.Name.Name)
	config := &types.Config{
		Context:                  nil,
		GoVersion:                "",
		IgnoreFuncBodies:         false,
		FakeImportC:              false,
		Error:                    nil,
		Importer:                 importer.Default(),
		Sizes:                    nil,
		DisableUnusedImportCheck: false,
	}

	ssaPackage, _, err := ssautil.BuildPackage(
		config, fileSet, pkg, files, ssa.SanityCheckFunctions,
	)
	if err != nil {
		return nil, fmt.Errorf("build SSA package: %w", err)
	}

	return ssaPackage, nil
}

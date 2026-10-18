package gossa

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Package struct {
	Go *ssa.Package
}

type Function struct {
	Go *ssa.Function
}

func (p Package) Functions() []Function {
	var functions []Function

	for _, m := range p.Go.Members {
		f, ok := m.(*ssa.Function)
		if !ok {
			continue
		}

		functions = append(functions, Function{f})
	}

	return functions
}

func (f Function) Entry() *ssa.BasicBlock {
	return f.Go.Blocks[0]
}

func ParseToSSA(text string) Package {
	path := "/dev/stdin"

	fileSet := token.NewFileSet()

	parseMode := parser.ParseComments

	file, err := parser.ParseFile(fileSet, path, text, parseMode)
	if err != nil {
		panic(fmt.Errorf("parse file: %w", err))
	}

	files := []*ast.File{file}
	pkg := types.NewPackage(path, file.Name.Name)
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

	buildMode := ssa.SanityCheckFunctions

	ssaPkg, _, err := ssautil.BuildPackage(config, fileSet, pkg, files, buildMode)
	if err != nil {
		panic(fmt.Errorf("build package: %w", err))
	}

	return Package{ssaPkg}
}

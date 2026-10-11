package gossa_test

import (
	"os"
	"testing"

	"github.com/vityaman-edu/gosvm/internal/gossa"
	"github.com/vityaman-edu/gosvm/internal/srcloc"
	"golang.org/x/tools/go/ssa"
)

func TestBuildSSA(t *testing.T) {
	t.Parallel()

	pkg := buildFixture(t)
	for _, name := range []string{
		"simpleIf", "ifElse", "nestedIf", "simpleLoop", "whileLoop",
		"nestedLoops", "loopWithBreakContinue", "switchExample",
		"multipleReturns", "withPanic", "factorial", "withGoto",
		"withDefer", "complexConditions", "main",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			checkFunction(t, pkg, name)
		})
	}
}

func buildFixture(t *testing.T) *ssa.Package {
	t.Helper()

	const path = "examples/test_functions.go"

	text, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read SSA fixture: %v", err)
	}

	pkg, err := gossa.BuildSSA(srcloc.Source{Path: path, Text: string(text)})
	if err != nil {
		t.Fatalf("build SSA: %v", err)
	}

	if pkg == nil {
		t.Fatal("BuildSSA returned a nil package")
	}

	if name := pkg.Pkg.Name(); name != "main" {
		t.Errorf("package name = %q, want main", name)
	}

	return pkg
}

func checkFunction(t *testing.T, pkg *ssa.Package, name string) {
	t.Helper()

	function := pkg.Func(name)
	if function == nil {
		t.Fatalf("SSA function %q is missing", name)
	}

	if len(function.Blocks) == 0 {
		t.Errorf("SSA function %q has no basic blocks", name)
	}
}

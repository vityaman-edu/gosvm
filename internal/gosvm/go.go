package gosvm

import (
	"errors"
	"fmt"
	"os"

	"github.com/vityaman-edu/gosvm/internal/gossa"
)

type Input struct {
	Text     string
	Function string
	PrintSSA bool
	PrintSMT bool
}

func Go(input Input) {
	pkg := gossa.ParseToSSA(input.Text)

	if input.PrintSSA {
		_, err := pkg.WriteTo(os.Stderr)
		if err != nil {
			panic(err)
		}
	}

	err := goSSA(pkg, input)
	if err != nil {
		if compact, ok := errors.AsType[*gossa.CompactError](err); ok {
			err = compact.Expand(pkg.FileSet)
		}

		panic(err)
	}
}

func goSSA(pkg gossa.Package, input Input) error {
	function, err := pkg.Function(input.Function)
	if err != nil {
		return fmt.Errorf("look up function %q: %w", input.Function, err)
	}

	err = loop(loopInput{
		Function: function,
		PrintSMT: input.PrintSMT,
	})
	if err != nil {
		return err
	}

	return nil
}

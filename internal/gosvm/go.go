package gosvm

import (
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
		if x, ok := err.(*gossa.CompactError); ok {
			err = x.Expand(pkg.FileSet)
		}

		panic(err)
	}
}

func goSSA(pkg gossa.Package, input Input) error {
	f, err := pkg.Function(input.Function)
	if err != nil {
		return err
	}

	err = loop(loopInput{
		Function: f,
		PrintSMT: input.PrintSMT,
	})
	if err != nil {
		return err
	}

	return nil
}

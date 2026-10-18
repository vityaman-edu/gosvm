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
}

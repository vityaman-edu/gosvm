package gossa

import (
	"fmt"
	"io"

	"golang.org/x/tools/go/ssa"
)

func (p Package) WriteTo(w io.Writer) (int64, error) {
	var n int64

	for _, m := range p.Go.Members {
		f, ok := m.(*ssa.Function)
		if !ok {
			continue
		}

		n1, err := Function{f}.WriteTo(w)
		if err != nil {
			return n, err
		}

		n += n1
	}

	return n, nil
}

func (f Function) WriteTo(w io.Writer) (int64, error) {
	n, err := f.Go.WriteTo(w)
	if err != nil {
		return n, fmt.Errorf("write SSA function: %w", err)
	}

	return n, nil
}

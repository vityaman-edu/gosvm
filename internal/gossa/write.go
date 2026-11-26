package gossa

import (
	"fmt"
	"io"

	"golang.org/x/tools/go/ssa"
)

func (p Package) WriteTo(writer io.Writer) (int64, error) {
	var size int64

	for _, m := range p.Go.Members {
		f, ok := m.(*ssa.Function)
		if !ok {
			continue
		}

		n, err := Function{f}.WriteTo(writer)
		if err != nil {
			return size, err
		}

		size += n
	}

	return size, nil
}

func (f Function) WriteTo(w io.Writer) (int64, error) {
	n, err := f.Go.WriteTo(w)
	if err != nil {
		return n, fmt.Errorf("write SSA function: %w", err)
	}

	return n, nil
}

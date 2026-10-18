package gossa

import (
	"fmt"
	"io"
)

func (p Package) WriteTo(w io.Writer) (int64, error) {
	n, err := p.Go.WriteTo(w)
	if err != nil {
		return n, fmt.Errorf("write SSA package: %w", err)
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

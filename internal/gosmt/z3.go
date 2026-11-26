package gosmt

import (
	"github.com/ebukreev/go-z3/z3"
	"golang.org/x/tools/go/ssa"
)

type Z3 struct {
	ctx    *z3.Context
	config *z3.Config
}

type Symbols = map[ssa.Value]z3.Value

func NewZ3() *Z3 {
	config := &z3.Config{}
	ctx := z3.NewContext(config)

	return &Z3{
		ctx:    ctx,
		config: config,
	}
}

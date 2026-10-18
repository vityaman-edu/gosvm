package gosmt

import "github.com/ebukreev/go-z3/z3"

type Z3 struct {
	Ctx    *z3.Context
	Config *z3.Config
	Vars   map[string]z3.Value
}

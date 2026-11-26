package gosmt

import "github.com/ebukreev/go-z3/z3"

func (z *Z3) Bool(b bool) z3.Bool {
	return z.ctx.FromBool(b)
}

func (z *Z3) Simplify(x z3.Value) z3.Value {
	return z.ctx.Simplify(x, z.config)
}

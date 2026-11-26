package gosmt

import "github.com/ebukreev/go-z3/z3"

func (z *Z3) Bool(b bool) z3.Value {
	return z.ctx.FromBool(b)
}

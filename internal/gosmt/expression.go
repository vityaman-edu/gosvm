package gosmt

import (
	"go/constant"
	"go/token"
	"go/types"

	"github.com/ebukreev/go-z3/z3"
	"github.com/vityaman-edu/gosvm/internal/gossa"
	"golang.org/x/tools/go/ssa"
)

func (z *Z3) TranslateValue(v ssa.Value, s Symbols) (z3.Value, error) {
	if value, ok := s[v]; ok {
		return value, nil
	}

	switch v := v.(type) {
	case *ssa.Alloc:
		return z.translateAlloc(v, s)
	case *ssa.BinOp:
		return z.translateBinOp(v, s)
	case *ssa.Builtin:
		return z.translateBuiltin(v, s)
	case *ssa.Call:
		return z.translateCall(v, s)
	case *ssa.ChangeInterface:
		return z.translateChangeInterface(v, s)
	case *ssa.ChangeType:
		return z.translateChangeType(v, s)
	case *ssa.Const:
		return z.translateConst(v, s)
	case *ssa.Convert:
		return z.translateConvert(v, s)
	case *ssa.Extract:
		return z.translateExtract(v, s)
	case *ssa.Field:
		return z.translateField(v, s)
	case *ssa.FieldAddr:
		return z.translateFieldAddr(v, s)
	case *ssa.FreeVar:
		return z.translateFreeVar(v, s)
	case *ssa.Function:
		return z.translateFunction(v, s)
	case *ssa.Global:
		return z.translateGlobal(v, s)
	case *ssa.Index:
		return z.translateIndex(v, s)
	case *ssa.IndexAddr:
		return z.translateIndexAddr(v, s)
	case *ssa.Lookup:
		return z.translateLookup(v, s)
	case *ssa.MakeChan:
		return z.translateMakeChan(v, s)
	case *ssa.MakeClosure:
		return z.translateMakeClosure(v, s)
	case *ssa.MakeInterface:
		return z.translateMakeInterface(v, s)
	case *ssa.MakeMap:
		return z.translateMakeMap(v, s)
	case *ssa.MakeSlice:
		return z.translateMakeSlice(v, s)
	case *ssa.MultiConvert:
		return z.translateMultiConvert(v, s)
	case *ssa.Next:
		return z.translateNext(v, s)
	case *ssa.Parameter:
		return z.translateParameter(v, s)
	case *ssa.Phi:
		return z.translatePhi(v, s)
	case *ssa.Range:
		return z.translateRange(v, s)
	case *ssa.Select:
		return z.translateSelect(v, s)
	case *ssa.Slice:
		return z.translateSlice(v, s)
	case *ssa.SliceToArrayPointer:
		return z.translateSliceToArrayPointer(v, s)
	case *ssa.TypeAssert:
		return z.translateTypeAssert(v, s)
	case *ssa.UnOp:
		return z.translateUnOp(v, s)
	default:
		return nil, gossa.Errorf("unsupported SSA value type %T", v.Pos(), v)
	}
}

func (z *Z3) translateAlloc(v *ssa.Alloc, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Alloc is not supported", v.Pos())
}

func (z *Z3) translateBinOp(v *ssa.BinOp, s Symbols) (z3.Value, error) {
	x, err := z.TranslateValue(v.X, s)
	if err != nil {
		return nil, err
	}

	y, err := z.TranslateValue(v.Y, s)
	if err != nil {
		return nil, err
	}

	switch v.Op {
	case token.ADD:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.Add(y), nil

	case token.SUB:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.Sub(y), nil

	case token.MUL:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.Mul(y), nil

	case token.QUO:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.Div(y), nil

	case token.REM:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.Mod(y), nil

	case token.AND:
		x := x.(z3.Bool)
		y := y.(z3.Bool)
		return x.And(y), nil

	case token.OR:
		x := x.(z3.Bool)
		y := y.(z3.Bool)
		return x.Or(y), nil

	case token.XOR:
		return nil, gossa.Errorf("BinOp %s is not supported", v.Pos(), v.Op)

	case token.SHL:
		return nil, gossa.Errorf("BinOp %s is not supported", v.Pos(), v.Op)

	case token.SHR:
		return nil, gossa.Errorf("BinOp %s is not supported", v.Pos(), v.Op)

	case token.AND_NOT:
		return nil, gossa.Errorf("BinOp %s is not supported", v.Pos(), v.Op)

	case token.EQL:
		switch x := x.(type) {
		case z3.Int:
			y := y.(z3.Int)
			return x.Eq(y), nil
		case z3.Bool:
			y := y.(z3.Bool)
			return x.Eq(y), nil
		default:
			return nil, gossa.Errorf("BinOp %s unexpected x type", v.Pos(), v.Op)
		}

	case token.NEQ:
		switch x := x.(type) {
		case z3.Int:
			y := y.(z3.Int)
			return x.NE(y), nil
		case z3.Bool:
			y := y.(z3.Bool)
			return x.NE(y), nil
		default:
			return nil, gossa.Errorf("BinOp %s unexpected x type", v.Pos(), v.Op)
		}

	case token.LSS:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.LT(y), nil

	case token.LEQ:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.LE(y), nil

	case token.GTR:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.GT(y), nil

	case token.GEQ:
		x := x.(z3.Int)
		y := y.(z3.Int)
		return x.GE(y), nil

	default:
		return nil, gossa.Errorf("unknown BinOp %s", v.Pos(), v.Op)
	}
}

func (z *Z3) translateBuiltin(v *ssa.Builtin, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Builtin is not supported", v.Pos())
}

func (z *Z3) translateCall(v *ssa.Call, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Call is not supported", v.Pos())
}

func (z *Z3) translateChangeInterface(v *ssa.ChangeInterface, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("ChangeInterface is not supported", v.Pos())
}

func (z *Z3) translateChangeType(v *ssa.ChangeType, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("ChangeType is not supported", v.Pos())
}

func (z *Z3) translateConst(v *ssa.Const, _ Symbols) (z3.Value, error) {
	switch v.Value.Kind() {
	case constant.Unknown:
		return nil, gossa.Errorf("unknown constant kind", v.Pos())

	case constant.Bool:
		x := constant.BoolVal(v.Value)
		return z.ctx.FromBool(x), nil

	case constant.String:
		return nil, gossa.Errorf("String constant is not supported", v.Pos())

	case constant.Int:
		x, ok := constant.Int64Val(v.Value)
		if !ok {
			return nil, gossa.Errorf("Too large Int constant", v.Pos())
		}

		return z.ctx.FromInt(x, z.ctx.IntSort()), nil

	case constant.Float:
		return nil, gossa.Errorf("Float constant is not supported", v.Pos())

	case constant.Complex:
		return nil, gossa.Errorf("Complex constant is not supported", v.Pos())

	default:
		return nil, gossa.Errorf("unsupported constant kind %d", v.Pos(), v.Value.Kind())
	}
}

func (z *Z3) translateConvert(v *ssa.Convert, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Convert is not supported", v.Pos())
}

func (z *Z3) translateExtract(v *ssa.Extract, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Extract is not supported", v.Pos())
}

func (z *Z3) translateField(v *ssa.Field, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Field is not supported", v.Pos())
}

func (z *Z3) translateFieldAddr(v *ssa.FieldAddr, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("FieldAddr is not supported", v.Pos())
}

func (z *Z3) translateFreeVar(v *ssa.FreeVar, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("FreeVar is not supported", v.Pos())
}

func (z *Z3) translateFunction(v *ssa.Function, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Function is not supported", v.Pos())
}

func (z *Z3) translateGlobal(v *ssa.Global, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Global %s is not supported", v.Pos(), v.Name())
}

func (z *Z3) translateIndex(v *ssa.Index, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Index is not supported", v.Pos())
}

func (z *Z3) translateIndexAddr(v *ssa.IndexAddr, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("IndexAddr is not supported", v.Pos())
}

func (z *Z3) translateLookup(v *ssa.Lookup, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Lookup is not supported", v.Pos())
}

func (z *Z3) translateMakeChan(v *ssa.MakeChan, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("MakeChan is not supported", v.Pos())
}

func (z *Z3) translateMakeClosure(v *ssa.MakeClosure, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("MakeClosure is not supported", v.Pos())
}

func (z *Z3) translateMakeInterface(v *ssa.MakeInterface, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("MakeInterface is not supported", v.Pos())
}

func (z *Z3) translateMakeMap(v *ssa.MakeMap, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("MakeMap is not supported", v.Pos())
}

func (z *Z3) translateMakeSlice(v *ssa.MakeSlice, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("MakeSlice is not supported", v.Pos())
}

func (z *Z3) translateMultiConvert(v *ssa.MultiConvert, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("MultiConvert is not supported", v.Pos())
}

func (z *Z3) translateNext(v *ssa.Next, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Next is not supported", v.Pos())
}

func (z *Z3) translateParameter(v *ssa.Parameter, s Symbols) (z3.Value, error) {
	basic, ok := v.Type().Underlying().(*types.Basic)
	if !ok {
		return nil, gossa.Errorf("Parameter type is not supported", v.Pos())
	}

	var sort z3.Sort
	switch basic.Kind() {
	case types.Int:
		sort = z.ctx.IntSort()
	case types.Bool:
		sort = z.ctx.BoolSort()
	default:
		return nil, gossa.Errorf("Parameter type is not supported", v.Pos())
	}

	x := z.ctx.Const(v.Name(), sort)
	s[v] = x
	return x, nil
}

func (z *Z3) translatePhi(v *ssa.Phi, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Phi is not supported", v.Pos())
}

func (z *Z3) translateRange(v *ssa.Range, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Range is not supported", v.Pos())
}

func (z *Z3) translateSelect(v *ssa.Select, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Select is not supported", v.Pos())
}

func (z *Z3) translateSlice(v *ssa.Slice, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Slice is not supported", v.Pos())
}

func (z *Z3) translateSliceToArrayPointer(v *ssa.SliceToArrayPointer, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("SliceToArrayPointer is not supported", v.Pos())
}

func (z *Z3) translateTypeAssert(v *ssa.TypeAssert, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("TypeAssert is not supported", v.Pos())
}

func (z *Z3) translateUnOp(v *ssa.UnOp, s Symbols) (z3.Value, error) {
	x, err := z.TranslateValue(v.X, s)
	if err != nil {
		return nil, err
	}

	switch v.Op {
	case token.NOT:
		x := x.(z3.Bool)
		return x.Not(), nil
	case token.SUB:
		x := x.(z3.Int)
		return x.Neg(), nil
	case token.ARROW:
		return nil, gossa.Errorf("ARROW UnOp is not supported", v.Pos())
	case token.MUL:
		return nil, gossa.Errorf("MUL UnOp is not supported", v.Pos())
	case token.XOR:
		return nil, gossa.Errorf("XOR UnOp is not supported", v.Pos())
	default:
		return nil, gossa.Errorf("unknown UnOp %s", v.Pos(), v.Op)
	}
}

package gosmt

import (
	"go/constant"
	"go/token"
	"go/types"

	"github.com/ebukreev/go-z3/z3"
	"github.com/vityaman-edu/gosvm/internal/gossa"
	"golang.org/x/tools/go/ssa"
)

//nolint:cyclop,gocyclo,funlen // Each SSA value type dispatches to its own handler.
func (z *Z3) TranslateValue(value ssa.Value, symbols Symbols) (z3.Value, error) {
	if cached, ok := symbols[value]; ok {
		return cached, nil
	}

	switch value := value.(type) {
	case *ssa.Alloc:
		return z.translateAlloc(value, symbols)
	case *ssa.BinOp:
		return z.translateBinOp(value, symbols)
	case *ssa.Builtin:
		return z.translateBuiltin(value, symbols)
	case *ssa.Call:
		return z.translateCall(value, symbols)
	case *ssa.ChangeInterface:
		return z.translateChangeInterface(value, symbols)
	case *ssa.ChangeType:
		return z.translateChangeType(value, symbols)
	case *ssa.Const:
		return z.translateConst(value, symbols)
	case *ssa.Convert:
		return z.translateConvert(value, symbols)
	case *ssa.Extract:
		return z.translateExtract(value, symbols)
	case *ssa.Field:
		return z.translateField(value, symbols)
	case *ssa.FieldAddr:
		return z.translateFieldAddr(value, symbols)
	case *ssa.FreeVar:
		return z.translateFreeVar(value, symbols)
	case *ssa.Function:
		return z.translateFunction(value, symbols)
	case *ssa.Global:
		return z.translateGlobal(value, symbols)
	case *ssa.Index:
		return z.translateIndex(value, symbols)
	case *ssa.IndexAddr:
		return z.translateIndexAddr(value, symbols)
	case *ssa.Lookup:
		return z.translateLookup(value, symbols)
	case *ssa.MakeChan:
		return z.translateMakeChan(value, symbols)
	case *ssa.MakeClosure:
		return z.translateMakeClosure(value, symbols)
	case *ssa.MakeInterface:
		return z.translateMakeInterface(value, symbols)
	case *ssa.MakeMap:
		return z.translateMakeMap(value, symbols)
	case *ssa.MakeSlice:
		return z.translateMakeSlice(value, symbols)
	case *ssa.MultiConvert:
		return z.translateMultiConvert(value, symbols)
	case *ssa.Next:
		return z.translateNext(value, symbols)
	case *ssa.Parameter:
		return z.translateParameter(value, symbols)
	case *ssa.Phi:
		return z.translatePhi(value, symbols)
	case *ssa.Range:
		return z.translateRange(value, symbols)
	case *ssa.Select:
		return z.translateSelect(value, symbols)
	case *ssa.Slice:
		return z.translateSlice(value, symbols)
	case *ssa.SliceToArrayPointer:
		return z.translateSliceToArrayPointer(value, symbols)
	case *ssa.TypeAssert:
		return z.translateTypeAssert(value, symbols)
	case *ssa.UnOp:
		return z.translateUnOp(value, symbols)
	default:
		return nil, gossa.Errorf("unsupported SSA value type %T", value.Pos(), value)
	}
}

func (z *Z3) translateAlloc(v *ssa.Alloc, _ Symbols) (z3.Value, error) {
	return nil, gossa.Errorf("Alloc is not supported", v.Pos())
}

func (z *Z3) translateBinOp(operation *ssa.BinOp, symbols Symbols) (z3.Value, error) {
	lhs, err := z.TranslateValue(operation.X, symbols)
	if err != nil {
		return nil, err
	}

	rhs, err := z.TranslateValue(operation.Y, symbols)
	if err != nil {
		return nil, err
	}

	var result z3.Value

	switch operation.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		result, err = applyIntArithmetic(operation, lhs, rhs)
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		result, err = applyComparison(operation, lhs, rhs)
	case token.AND, token.OR:
		result, err = applyBoolBinOp(operation, lhs, rhs)
	case token.XOR, token.SHL, token.SHR, token.AND_NOT:
		return nil, gossa.Errorf("BinOp %s is not supported", operation.Pos(), operation.Op)
	default:
		return nil, gossa.Errorf("unknown BinOp %s", operation.Pos(), operation.Op)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

type intOperands struct {
	lhs, rhs z3.Int
}

func checkedIntOperands(operation *ssa.BinOp, lhs, rhs z3.Value) (intOperands, error) {
	intLHS, lhsOK := lhs.(z3.Int)
	intRHS, rhsOK := rhs.(z3.Int)

	if !lhsOK || !rhsOK {
		return intOperands{}, gossa.Errorf(
			"BinOp %s expects integer operands, got %T and %T", operation.Pos(), operation.Op, lhs, rhs,
		)
	}

	return intOperands{lhs: intLHS, rhs: intRHS}, nil
}

func applyIntArithmetic(operation *ssa.BinOp, lhs, rhs z3.Value) (z3.Int, error) {
	operands, err := checkedIntOperands(operation, lhs, rhs)
	if err != nil {
		return z3.Int{}, err
	}

	switch operation.Op {
	case token.ADD:
		return operands.lhs.Add(operands.rhs), nil
	case token.SUB:
		return operands.lhs.Sub(operands.rhs), nil
	case token.MUL:
		return operands.lhs.Mul(operands.rhs), nil
	case token.QUO:
		return operands.lhs.Div(operands.rhs), nil
	case token.REM:
		return operands.lhs.Mod(operands.rhs), nil
	default:
		panic("unexpected integer arithmetic operator")
	}
}

func applyComparison(operation *ssa.BinOp, lhs, rhs z3.Value) (z3.Bool, error) {
	if operation.Op == token.EQL || operation.Op == token.NEQ {
		if _, ok := lhs.(z3.Bool); ok {
			return applyBoolBinOp(operation, lhs, rhs)
		}
	}

	return applyIntComparison(operation, lhs, rhs)
}

func applyIntComparison(operation *ssa.BinOp, lhs, rhs z3.Value) (z3.Bool, error) {
	operands, err := checkedIntOperands(operation, lhs, rhs)
	if err != nil {
		return z3.Bool{}, err
	}

	switch operation.Op {
	case token.EQL:
		return operands.lhs.Eq(operands.rhs), nil
	case token.NEQ:
		return operands.lhs.NE(operands.rhs), nil
	case token.LSS:
		return operands.lhs.LT(operands.rhs), nil
	case token.LEQ:
		return operands.lhs.LE(operands.rhs), nil
	case token.GTR:
		return operands.lhs.GT(operands.rhs), nil
	case token.GEQ:
		return operands.lhs.GE(operands.rhs), nil
	default:
		panic("unexpected integer comparison operator")
	}
}

func applyBoolBinOp(operation *ssa.BinOp, lhs, rhs z3.Value) (z3.Bool, error) {
	boolLHS, lhsOK := lhs.(z3.Bool)
	boolRHS, rhsOK := rhs.(z3.Bool)

	if !lhsOK || !rhsOK {
		return z3.Bool{}, gossa.Errorf(
			"BinOp %s expects boolean operands, got %T and %T", operation.Pos(), operation.Op, lhs, rhs,
		)
	}

	switch operation.Op {
	case token.AND:
		return boolLHS.And(boolRHS), nil
	case token.OR:
		return boolLHS.Or(boolRHS), nil
	case token.EQL:
		return boolLHS.Eq(boolRHS), nil
	case token.NEQ:
		return boolLHS.NE(boolRHS), nil
	default:
		panic("unexpected boolean operator")
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

func (z *Z3) translateConst(literal *ssa.Const, _ Symbols) (z3.Value, error) {
	switch literal.Value.Kind() {
	case constant.Unknown:
		return nil, gossa.Errorf("unknown constant kind", literal.Pos())

	case constant.Bool:
		x := constant.BoolVal(literal.Value)

		return z.ctx.FromBool(x), nil

	case constant.String:
		return nil, gossa.Errorf("String constant is not supported", literal.Pos())

	case constant.Int:
		x, ok := constant.Int64Val(literal.Value)
		if !ok {
			return nil, gossa.Errorf("Too large Int constant", literal.Pos())
		}

		return z.ctx.FromInt(x, z.ctx.IntSort()), nil

	case constant.Float:
		return nil, gossa.Errorf("Float constant is not supported", literal.Pos())

	case constant.Complex:
		return nil, gossa.Errorf("Complex constant is not supported", literal.Pos())

	default:
		return nil, gossa.Errorf("unsupported constant kind %d", literal.Pos(), literal.Value.Kind())
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

func (z *Z3) translateParameter(parameter *ssa.Parameter, symbols Symbols) (z3.Value, error) {
	basic, ok := parameter.Type().Underlying().(*types.Basic)
	if !ok {
		return nil, gossa.Errorf("Parameter type is not supported", parameter.Pos())
	}

	var sort z3.Sort

	switch basic.Kind() {
	case types.Int:
		sort = z.ctx.IntSort()
	case types.Bool:
		sort = z.ctx.BoolSort()
	default:
		return nil, gossa.Errorf("Parameter type is not supported", parameter.Pos())
	}

	x := z.ctx.Const(parameter.Name(), sort)
	symbols[parameter] = x

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

func (z *Z3) translateUnOp(operation *ssa.UnOp, symbols Symbols) (z3.Value, error) {
	operand, err := z.TranslateValue(operation.X, symbols)
	if err != nil {
		return nil, err
	}

	switch operation.Op {
	case token.NOT:
		boolean, ok := operand.(z3.Bool)
		if !ok {
			return nil, gossa.Errorf("UnOp %s expects a boolean operand, got %T", operation.Pos(), operation.Op, operand)
		}

		return boolean.Not(), nil
	case token.SUB:
		integer, ok := operand.(z3.Int)
		if !ok {
			return nil, gossa.Errorf("UnOp %s expects an integer operand, got %T", operation.Pos(), operation.Op, operand)
		}

		return integer.Neg(), nil
	case token.ARROW:
		return nil, gossa.Errorf("ARROW UnOp is not supported", operation.Pos())
	case token.MUL:
		return nil, gossa.Errorf("MUL UnOp is not supported", operation.Pos())
	case token.XOR:
		return nil, gossa.Errorf("XOR UnOp is not supported", operation.Pos())
	default:
		return nil, gossa.Errorf("unknown UnOp %s", operation.Pos(), operation.Op)
	}
}

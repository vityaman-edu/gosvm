package translation

import (
	"fmt"

	"github.com/ebukreev/go-z3/z3"
	"github.com/vityaman-edu/gosvm/internal/fp"
	"github.com/vityaman-edu/gosvm/internal/symbolic"
)

type Z3Translator struct {
	ctx    *z3.Context
	config *z3.Config
	vars   map[string]z3.Value
}

func NewZ3Translator() *Z3Translator {
	config := &z3.Config{}
	ctx := z3.NewContext(config)

	return &Z3Translator{
		ctx:    ctx,
		config: config,
		vars:   make(map[string]z3.Value),
	}
}

var _ symbolic.ExpressionVisitor = (*Z3Translator)(nil)

func (t *Z3Translator) VisitVariable(variable *symbolic.Variable) any {
	if value, ok := t.vars[variable.Name]; ok {
		return value
	}

	factory := &VariableTypeVisitor{t.ctx}

	constructor, ok := variable.Type.Accept(factory).(func(string) z3.Value)
	if !ok {
		panic("variable type visitor must return a Z3 variable constructor")
	}

	t.vars[variable.Name] = constructor(variable.Name)

	return t.VisitVariable(variable)
}

func (t *Z3Translator) VisitIntConstant(constant *symbolic.IntConstant) any {
	return t.ctx.FromInt(constant.Value, t.ctx.IntSort())
}

func (t *Z3Translator) VisitBoolConstant(constant *symbolic.BoolConstant) any {
	return t.ctx.FromBool(constant.Value)
}

func (t *Z3Translator) VisitBinaryOperation(operation *symbolic.BinaryOperation) any {
	lhs, isValue := operation.LHS.Accept(t).(z3.Value)
	if !isValue {
		panic("left operand visitor must return a Z3 value")
	}

	rhs, isValue := operation.RHS.Accept(t).(z3.Value)
	if !isValue {
		panic("right operand visitor must return a Z3 value")
	}

	if lhsInt, isInt := lhs.(z3.Int); isInt {
		if rhsInt, isInt := rhs.(z3.Int); isInt {
			switch operation.Operator {
			case symbolic.ADD, symbolic.SUB, symbolic.MUL, symbolic.DIV, symbolic.MOD:
				return translateIntArithmetic(operation.Operator, lhsInt, rhsInt)
			case symbolic.LT, symbolic.LE, symbolic.GT, symbolic.GE, symbolic.EQ, symbolic.NE:
				return translateIntComparison(operation.Operator, lhsInt, rhsInt)
			default:
				panic(fmt.Sprintf("unexpected binary operator %d", operation.Operator))
			}
		}
	}

	if lhsBool, isBool := lhs.(z3.Bool); isBool {
		if rhsBool, isBool := rhs.(z3.Bool); isBool {
			return translateBoolEquality(operation.Operator, lhsBool, rhsBool)
		}
	}

	panic(fmt.Sprintf(
		"binary operator %s: unsupported operand types (%T, %T); "+
			"expected two integers, or two booleans for equality/inequality",
		operation.Operator.String(), lhs, rhs,
	))
}

func translateIntArithmetic(operator symbolic.BinaryOperator, lhs, rhs z3.Int) z3.Int {
	switch operator {
	case symbolic.ADD:
		return lhs.Add(rhs)
	case symbolic.SUB:
		return lhs.Sub(rhs)
	case symbolic.MUL:
		return lhs.Mul(rhs)
	case symbolic.DIV:
		return lhs.Div(rhs)
	case symbolic.MOD:
		return lhs.Mod(rhs)
	case symbolic.LT, symbolic.LE, symbolic.GT, symbolic.GE, symbolic.EQ, symbolic.NE:
	}

	panic(fmt.Sprintf("unexpected binary operator %d", operator))
}

func translateIntComparison(operator symbolic.BinaryOperator, lhs, rhs z3.Int) z3.Bool {
	switch operator {
	case symbolic.LT:
		return lhs.LT(rhs)
	case symbolic.LE:
		return lhs.LE(rhs)
	case symbolic.GT:
		return lhs.GT(rhs)
	case symbolic.GE:
		return lhs.GE(rhs)
	case symbolic.EQ:
		return lhs.Eq(rhs)
	case symbolic.NE:
		return lhs.NE(rhs)
	case symbolic.ADD, symbolic.SUB, symbolic.MUL, symbolic.DIV, symbolic.MOD:
	}

	panic(fmt.Sprintf("unexpected binary operator %d", operator))
}

func translateBoolEquality(operator symbolic.BinaryOperator, lhs, rhs z3.Bool) z3.Bool {
	switch operator {
	case symbolic.EQ:
		return lhs.Eq(rhs)
	case symbolic.NE:
		return lhs.NE(rhs)
	case symbolic.ADD, symbolic.SUB, symbolic.MUL, symbolic.DIV, symbolic.MOD,
		symbolic.LT, symbolic.LE, symbolic.GT, symbolic.GE:
	}

	panic(fmt.Sprintf(
		"binary operator %s: unsupported operand types (%T, %T); "+
			"expected two integers, or two booleans for equality/inequality",
		operator.String(), lhs, rhs,
	))
}

func (t *Z3Translator) VisitLogicalOperation(operation *symbolic.LogicalOperation) any {
	const implicationOperandCount = 2

	operands := fp.Map(operation.Operands, func(operand symbolic.Expression) z3.Bool {
		value, ok := operand.Accept(t).(z3.Bool)
		if !ok {
			panic("logical operand visitor must return a Z3 boolean")
		}

		return value
	})

	switch operation.Operator {
	case symbolic.AND:
		var result z3.Bool
		fp.Fold1(operands, &result, func(lhs, rhs z3.Bool) z3.Bool {
			return lhs.And(rhs)
		})

		return result
	case symbolic.OR:
		var result z3.Bool
		fp.Fold1(operands, &result, func(lhs, rhs z3.Bool) z3.Bool {
			return lhs.Or(rhs)
		})

		return result
	case symbolic.NOT:
		if len(operands) != 1 {
			panic("NOT requires exactly one operand")
		}

		return operands[0].Not()
	case symbolic.IMPLIES:
		if len(operands) != implicationOperandCount {
			panic("IMPLIES requires exactly two operands")
		}

		return operands[0].Implies(operands[1])
	default:
		panic(fmt.Sprintf("unexpected logical operator %d", operation.Operator))
	}
}

type VariableTypeVisitor struct {
	ctx *z3.Context
}

var _ symbolic.TypeVisitor = (*VariableTypeVisitor)(nil)

func (t *VariableTypeVisitor) VisitIntType(_ *symbolic.IntType) any {
	return func(name string) z3.Value {
		return t.ctx.IntConst(name)
	}
}

func (t *VariableTypeVisitor) VisitBoolType(_ *symbolic.BoolType) any {
	return func(name string) z3.Value {
		return t.ctx.BoolConst(name)
	}
}

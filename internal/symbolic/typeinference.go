package symbolic

import (
	"errors"
	"fmt"
)

var ErrNilExpression = errors.New("cannot infer the type of a nil expression")

var (
	errUnknownType            = errors.New("unknown type")
	errUnknownBinaryOperator  = errors.New("unknown binary operator")
	errIncompatibleOperands   = errors.New("incompatible operand types")
	errNegationArity          = errors.New("NOT requires exactly one operand")
	errImplicationArity       = errors.New("IMPLIES requires exactly two operands")
	errUnknownLogicalOperator = errors.New("unknown logical operator")
	errEmptyLogicalOperation  = errors.New("logical operation requires at least one operand")
	errExpectedBool           = errors.New("must be bool")
)

func TypeInfer(expression Expression) TypeInferenceResult {
	if expression == nil {
		return TypeInferenceResult{Type: nil, Err: ErrNilExpression}
	}

	result, ok := expression.Accept(TypeInferenceVisitor{}).(TypeInferenceResult)
	if !ok {
		panic("type inference visitor must return a type inference result")
	}

	return result
}

type TypeInferenceResult struct {
	Type Type
	Err  error
}

type TypeInferenceVisitor struct{}

var _ ExpressionVisitor = TypeInferenceVisitor{}

func (TypeInferenceVisitor) VisitVariable(variable *Variable) any {
	switch variable.Type.(type) {
	case *IntType, *BoolType:
		return TypeInferenceResult{Type: variable.Type, Err: nil}
	default:
		return TypeInferenceResult{Type: nil, Err: fmt.Errorf("variable %q has an %w", variable.Name, errUnknownType)}
	}
}

func (TypeInferenceVisitor) VisitIntConstant(_ *IntConstant) any {
	return TypeInferenceResult{Type: &IntType{}, Err: nil}
}

func (TypeInferenceVisitor) VisitBoolConstant(_ *BoolConstant) any {
	return TypeInferenceResult{Type: &BoolType{}, Err: nil}
}

func (TypeInferenceVisitor) VisitBinaryOperation(operation *BinaryOperation) any {
	lhs := TypeInfer(operation.LHS)
	if lhs.Err != nil {
		return lhs
	}

	rhs := TypeInfer(operation.RHS)
	if rhs.Err != nil {
		return rhs
	}

	return inferBinaryOperationType(operation.Operator, lhs.Type, rhs.Type)
}

func inferBinaryOperationType(operator BinaryOperator, lhs, rhs Type) TypeInferenceResult {
	_, lhsInt := lhs.(*IntType)
	_, rhsInt := rhs.(*IntType)
	bothInt := lhsInt && rhsInt

	switch operator {
	case ADD, SUB, MUL, DIV, MOD:
		if bothInt {
			return TypeInferenceResult{Type: &IntType{}, Err: nil}
		}
	case LT, LE, GT, GE:
		if bothInt {
			return TypeInferenceResult{Type: &BoolType{}, Err: nil}
		}
	case EQ, NE:
		if haveSameOperandType(lhs, rhs) {
			return TypeInferenceResult{Type: &BoolType{}, Err: nil}
		}
	default:
		return TypeInferenceResult{Type: nil, Err: fmt.Errorf("%w %d", errUnknownBinaryOperator, operator)}
	}

	return TypeInferenceResult{
		Type: nil,
		Err:  fmt.Errorf("%w %T and %T for binary operator %d", errIncompatibleOperands, lhs, rhs, operator),
	}
}

func haveSameOperandType(lhs, rhs Type) bool {
	switch lhs.(type) {
	case *IntType:
		_, matches := rhs.(*IntType)

		return matches
	case *BoolType:
		_, matches := rhs.(*BoolType)

		return matches
	default:
		return false
	}
}

func (TypeInferenceVisitor) VisitLogicalOperation(operation *LogicalOperation) any {
	err := validateLogicalOperation(operation)
	if err != nil {
		return TypeInferenceResult{Type: nil, Err: err}
	}

	for index, operand := range operation.Operands {
		result := TypeInfer(operand)
		if result.Err != nil {
			return result
		}

		if _, ok := result.Type.(*BoolType); !ok {
			return TypeInferenceResult{
				Type: nil,
				Err:  fmt.Errorf("logical operand %d %w, got %T", index, errExpectedBool, result.Type),
			}
		}
	}

	return TypeInferenceResult{Type: &BoolType{}, Err: nil}
}

func validateLogicalOperation(operation *LogicalOperation) error {
	const implicationOperandCount = 2

	switch operation.Operator {
	case AND, OR:
	case NOT:
		if len(operation.Operands) != 1 {
			return errNegationArity
		}
	case IMPLIES:
		if len(operation.Operands) != implicationOperandCount {
			return errImplicationArity
		}
	default:
		return fmt.Errorf("%w %d", errUnknownLogicalOperator, operation.Operator)
	}

	if len(operation.Operands) == 0 {
		return errEmptyLogicalOperation
	}

	return nil
}

package symbolic

import (
	"strconv"
	"strings"
)

func ToStringExpression(expression Expression) string {
	result, ok := expression.Accept(ToStringExpressionVisitor{}).(string)
	if !ok {
		panic("ToStringExpression: visitor result is not a string")
	}

	return result
}

func ToStringType(t Type) string {
	result, ok := t.Accept(ToStringTypeVisitor{}).(string)
	if !ok {
		panic("ToStringType: visitor result is not a string")
	}

	return result
}

type ToStringExpressionVisitor struct{}

var _ ExpressionVisitor = ToStringExpressionVisitor{}

func (ToStringExpressionVisitor) VisitVariable(e *Variable) any {
	return e.Name
}

func (ToStringExpressionVisitor) VisitIntConstant(e *IntConstant) any {
	return strconv.FormatInt(e.Value, 10)
}

func (ToStringExpressionVisitor) VisitBoolConstant(e *BoolConstant) any {
	return strconv.FormatBool(e.Value)
}

func (operator BinaryOperator) String() string {
	operators := [...]string{
		ADD: "+",
		SUB: "-",
		MUL: "*",
		DIV: "/",
		MOD: "%",
		EQ:  "==",
		NE:  "!=",
		LT:  "<",
		LE:  "<=",
		GT:  ">",
		GE:  ">=",
	}

	return operators[operator]
}

func (operator LogicalOperator) String() string {
	operators := [...]string{
		AND:     "&&",
		OR:      "||",
		NOT:     "!",
		IMPLIES: "=>",
	}

	return operators[operator]
}

func (ToStringExpressionVisitor) VisitBinaryOperation(expression *BinaryOperation) any {
	return "(" + ToStringExpression(expression.LHS) + " " + expression.Operator.String() +
		" " + ToStringExpression(expression.RHS) + ")"
}

func (ToStringExpressionVisitor) VisitLogicalOperation(expression *LogicalOperation) any {
	if expression.Operator == NOT {
		return "(" + expression.Operator.String() + ToStringExpression(expression.Operands[0]) + ")"
	}

	operands := make([]string, len(expression.Operands))

	for i, operand := range expression.Operands {
		operands[i] = ToStringExpression(operand)
	}

	return "(" + strings.Join(operands, " "+expression.Operator.String()+" ") + ")"
}

type ToStringTypeVisitor struct{}

var _ TypeVisitor = ToStringTypeVisitor{}

func (ToStringTypeVisitor) VisitIntType(_ *IntType) any {
	return "int"
}

func (ToStringTypeVisitor) VisitBoolType(_ *BoolType) any {
	return "bool"
}

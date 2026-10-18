package symbolic

type Expression interface {
	Accept(v ExpressionVisitor) any
}

type ExpressionVisitor interface {
	VisitVariable(e *Variable) any

	VisitIntConstant(e *IntConstant) any

	VisitBoolConstant(e *BoolConstant) any

	VisitBinaryOperation(e *BinaryOperation) any

	VisitLogicalOperation(e *LogicalOperation) any
}

type Variable struct {
	Type Type
	Name string
}

func (e *Variable) Accept(v ExpressionVisitor) any {
	return v.VisitVariable(e)
}

type IntConstant struct {
	Value int64
}

func (e *IntConstant) Accept(v ExpressionVisitor) any {
	return v.VisitIntConstant(e)
}

type BoolConstant struct {
	Value bool
}

func (e *BoolConstant) Accept(v ExpressionVisitor) any {
	return v.VisitBoolConstant(e)
}

type BinaryOperation struct {
	LHS      Expression
	RHS      Expression
	Operator BinaryOperator
}

func (e *BinaryOperation) Accept(v ExpressionVisitor) any {
	return v.VisitBinaryOperation(e)
}

type LogicalOperation struct {
	Operands []Expression
	Operator LogicalOperator
}

func (e *LogicalOperation) Accept(v ExpressionVisitor) any {
	return v.VisitLogicalOperation(e)
}

type BinaryOperator int

const (
	ADD BinaryOperator = iota
	SUB
	MUL
	DIV
	MOD

	EQ
	NE
	LT
	LE
	GT
	GE
)

type LogicalOperator int

const (
	AND LogicalOperator = iota
	OR
	NOT
	IMPLIES
)

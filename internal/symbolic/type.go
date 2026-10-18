package symbolic

type Type interface {
	Accept(v TypeVisitor) any
}

type TypeVisitor interface {
	VisitIntType(t *IntType) any

	VisitBoolType(t *BoolType) any
}

type IntType struct{}

func (t *IntType) Accept(v TypeVisitor) any {
	return v.VisitIntType(t)
}

type BoolType struct{}

func (t *BoolType) Accept(v TypeVisitor) any {
	return v.VisitBoolType(t)
}

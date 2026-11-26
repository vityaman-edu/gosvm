package gosched

import "math"

type Eisenhower[T any] interface {
	Evaluate(value T) int
}

type DfsEisenhower[T any] struct {
	counter int
}

func NewDfsEisenhower[T any]() *DfsEisenhower[T] {
	return &DfsEisenhower[T]{counter: math.MinInt}
}

func (eisenhower *DfsEisenhower[T]) Evaluate(_ T) int {
	eisenhower.counter++

	return eisenhower.counter
}

type BfsEisenhower[T any] struct {
	counter int
}

func NewBfsEisenhower[T any]() *BfsEisenhower[T] {
	return &BfsEisenhower[T]{counter: math.MaxInt}
}

func (eisenhower *BfsEisenhower[T]) Evaluate(_ T) int {
	eisenhower.counter--

	return eisenhower.counter
}

package gosched

import "container/heap"

type Scheduler[T any] interface {
	Submit(value T)
	Poll() (T, bool)
}

type EisenhowerScheduler[T any] struct {
	priority Eisenhower[T]
	queue    priorityQueue[T]
}

func NewEisenhowerScheduler[T any](priority Eisenhower[T]) *EisenhowerScheduler[T] {
	return &EisenhowerScheduler[T]{priority: priority, queue: nil}
}

func (sched *EisenhowerScheduler[T]) Submit(value T) {
	heap.Push(&sched.queue, priorityItem[T]{
		value:    value,
		priority: sched.priority.Evaluate(value),
	})
}

//nolint:ireturn // Poll returns a value of the caller's chosen type T.
func (sched *EisenhowerScheduler[T]) Poll() (T, bool) {
	if sched.queue.Len() == 0 {
		var zero T

		return zero, false
	}

	//nolint:errcheck,forcetypeassert // Submit only inserts priorityItem[T] values.
	item := heap.Pop(&sched.queue).(priorityItem[T])

	return item.value, true
}

type LimitedScheduler[T any] struct {
	origin Scheduler[T]
	limit  int
	polled int
}

// NewLimitedScheduler limits calls to Poll, including polls of an empty scheduler.
func NewLimitedScheduler[T any](limit int, origin Scheduler[T]) *LimitedScheduler[T] {
	return &LimitedScheduler[T]{origin: origin, limit: limit, polled: 0}
}

func (sched *LimitedScheduler[T]) Submit(value T) {
	sched.origin.Submit(value)
}

//nolint:ireturn // Poll returns a value of the caller's chosen type T.
func (sched *LimitedScheduler[T]) Poll() (T, bool) {
	if sched.polled >= sched.limit {
		var zero T

		return zero, false
	}

	sched.polled++

	return sched.origin.Poll()
}

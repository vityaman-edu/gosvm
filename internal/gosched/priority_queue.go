package gosched

type priorityItem[T any] struct {
	value    T
	priority int
}

type priorityQueue[T any] []priorityItem[T]

func (pq priorityQueue[T]) Len() int { return len(pq) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue[T]) Push(x any) {
	//nolint:errcheck,forcetypeassert // Submit only inserts priorityItem[T] values.
	item := x.(priorityItem[T])

	*pq = append(*pq, item)
}

func (pq *priorityQueue[T]) Pop() any {
	items := *pq
	last := len(items) - 1
	item := items[last]
	clear(items[last:])
	*pq = items[:last]

	return item
}

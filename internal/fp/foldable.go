package fp

func Fold[T, U any](values []T, accumulator *U, reduce func(U, T) U) {
	result := *accumulator

	for _, value := range values {
		result = reduce(result, value)
	}

	*accumulator = result
}

func Fold1[T any](values []T, accumulator *T, reduce func(T, T) T) {
	if len(values) == 0 {
		panic("Fold1: empty slice")
	}

	result := values[0]
	Fold(values[1:], &result, reduce)
	*accumulator = result
}

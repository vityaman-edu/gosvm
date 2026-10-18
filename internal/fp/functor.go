package fp

func Map[T, U any](values []T, f func(T) U) []U {
	result := make([]U, len(values))

	for i, value := range values {
		result[i] = f(value)
	}

	return result
}

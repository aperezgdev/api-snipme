package pkg

func Map[T, Y any](array []T, f func(T) Y) []Y {
	var result []Y
	for _, item := range array {
		result = append(result, f(item))
	}
	return result
}

func GroupBy[T any, K comparable](array []T, f func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range array {
		key := f(item)
		result[key] = append(result[key], item)
	}
	return result
}

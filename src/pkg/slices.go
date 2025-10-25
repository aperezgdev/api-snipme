package pkg

func Map[T, Y any](array []T, f func(T) Y) []Y {
	var result []Y
	for _, item := range array {
		result = append(result, f(item))
	}
	return result
}

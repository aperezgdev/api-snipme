package pkg

type Optional[T any] struct {
	value   T
	present bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{present: false}
}

func (o Optional[T]) IsPresent() bool {
	return o.present
}

func (o Optional[T]) Get() T {
	if !o.present {
		panic("Optional is empty")
	}
	return o.value
}

func (o Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

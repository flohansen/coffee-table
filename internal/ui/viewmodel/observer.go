package viewmodel

type Listener[T any] func(value T)

type Observer[T any] struct {
	value     T
	listeners []Listener[T]
}

func NewObserver[T any](value T) *Observer[T] {
	return &Observer[T]{
		value: value,
	}
}

func (o *Observer[T]) Bind(lis func(T)) {
	o.listeners = append(o.listeners, lis)
	lis(o.value)
}

func (o *Observer[T]) Set(value T) {
	o.value = value
	for _, lis := range o.listeners {
		lis(value)
	}
}

func (o *Observer[T]) Get() T {
	return o.value
}

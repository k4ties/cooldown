package cooldown

// ValuedOption ...
type ValuedOption[T any] = func(*Valued[T])

// OptionHandler is used to set custom handlers to the Valued cooldown.
func OptionHandler[T any](h ValuedHandler[T]) ValuedOption[T] {
	return func(cd *Valued[T]) {
		cd.handler.Store(h)
	}
}

package cooldown

type (
	// ValuedOption is option implementation for the Valued cooldown.
	ValuedOption[T any] = func(cd *Valued[T])
	// Option is the option implementation for default CoolDown.
	Option = func(cd *CoolDown)
)

func ValuedOptionHandler[T any](h ValuedHandler[T]) ValuedOption[T] {
	return func(cd *Valued[T]) {
		cd.Handle(h)
	}
}

func OptionHandler(h Handler) Option {
	return func(cd *CoolDown) {
		cd.Handle(h)
	}
}

package cooldown

import "github.com/k4ties/cooldown/internal/event"

type (
	// ValuedContext is context for Valued.
	ValuedContext[T any] = event.Context[*Valued[T]]
	// ValuedOption is option implementation for the Valued cooldown.
	ValuedOption[T any] = func(cd *Valued[T])
)

type (
	// Context is context for CoolDown.
	Context = event.Context[*CoolDown]
	// Option is the option implementation for default CoolDown.
	Option = func(cd *CoolDown)
)

// convertContext creates new ValuedContext instance with T and provided valued
// cooldown.
func createContext[T any](cooldown *Valued[T]) *ValuedContext[T] {
	raw := event.C(cooldown)
	return raw
}

/*
  Options implementations
*/

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

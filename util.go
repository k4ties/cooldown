package cooldown

import "github.com/k4ties/cooldown/internal/event"

// convertContext creates new *event.Context with T. ValuedContext is
// technically equal for this type, but for go compiler it is hard to
// understand, so we need to additionally convert it to ValuedContext.
func createContext[T any](cooldown *Valued[T]) *ValuedContext[T] {
	raw := event.C(cooldown)
	ctx := (*ValuedContext[T])(
		raw,
	)
	return ctx
}

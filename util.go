package cooldown

import "github.com/k4ties/cooldown/internal/event"

// convertContext ...
func convertContext[T any](raw *event.Context[*Valued[T]]) *ValuedContext[T] {
	return (*ValuedContext[T])(raw)
}

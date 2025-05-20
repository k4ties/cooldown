package cooldown

import (
	"github.com/k4ties/cooldown/internal/event"
	"time"
)

// TicksPerSecond is amount how many times ticker will tick in a second.
const TicksPerSecond = 20

// tickDuration ...
func tickDuration() time.Duration {
	return time.Second / TicksPerSecond
}

// convertContext ...
func convertContext[T any](raw *event.Context[*Valued[T]]) *ValuedContext[T] {
	return (*ValuedContext[T])(raw)
}

package cooldown

import (
	"github.com/k4ties/cooldown/internal/atomic"
	"github.com/k4ties/cooldown/internal/event"
	"time"
)

// Valued represents basic cooldown with renew ability and custom values.
//
// Point of values is to use it in transactions and allow other users to handle the context value
// easily. All actions of the cooldown (like stop, start, renew) accepts a value. My personal
// opinion is to use in dragonfly (github.com/df-mc/dragonfly) world transactions. If you don't
// want to use values, there is zero implementation: CoolDown.
type Valued[T any] struct {
	basic *Basic

	duration atomic.Value[time.Duration]
	handler  atomic.Value[ValuedHandler[T]]
}

// NewValued creates new Valued cooldown.
func NewValued[T any](opts ...ValuedOption[T]) *Valued[T] {
	cooldown := &Valued[T]{}
	cooldown.basic = NewBasic()
	cooldown.duration = atomic.NewValue[time.Duration]()
	cooldown.handler = atomic.NewValue[ValuedHandler[T]]()
	for _, opt := range opts {
		opt(cooldown)
	}
	return cooldown
}

// Renew renews the cooldown, if it is currently active.
func (cooldown *Valued[T]) Renew(val T) {
	if !cooldown.Active() {
		// Not active, cannot renew
		return
	}

	duration, ok := cooldown.duration.Load()
	if duration == -1 || !ok {
		// Failed to load duration
		return
	}

	ctx := convertContext(event.C(cooldown))
	if cooldown.Handler().HandleRenew(ctx, val); ctx.Cancelled() {
		return
	}

	cooldown.basic.Set(duration)
}

// Start starts the cooldown, if it is not currently active.
func (cooldown *Valued[T]) Start(dur time.Duration, val T) {
	if cooldown.Active() {
		// Already active, cannot start again
		return
	}

	ctx := convertContext(event.C(cooldown))
	if cooldown.Handler().HandleStart(ctx, val); ctx.Cancelled() {
		return
	}

	cooldown.duration.Store(dur)
	cooldown.basic.Set(dur)

	Proc.Append(cooldown)
}

// Stop stops the cooldown, if it is currently active.
func (cooldown *Valued[T]) Stop(val T) {
	if !cooldown.Active() {
		// Not active, cannot stop
		return
	}

	cause := StopCauseCancelled
	cooldown.Handler().HandleStop(cooldown, StopCauseCancelled, val)

	Proc.Remove(cooldown)
	cooldown.stop(cause, false) // already handled
}

// stop implements processable.
func (cooldown *Valued[T]) stop(cause StopCause, handle bool) {
	if handle {
		var zeroT T
		cooldown.Handler().HandleStop(cooldown, cause, zeroT)
	}

	cooldown.duration.Store(-1)
	cooldown.basic.Reset()
}

// Handler returns current cooldown handler. If it is not set, NopHandler will be returned.
func (cooldown *Valued[T]) Handler() ValuedHandler[T] {
	handler, ok := cooldown.handler.Load()
	if !ok || handler == nil {
		// Should never happen anyway
		handler = NopValuedHandler[T]{}
	}
	return handler
}

// Handle updates current cooldown handler. If user entered nil as argument, current handler will
// be removed.
func (cooldown *Valued[T]) Handle(handler ValuedHandler[T]) {
	// Making sure h is never nil
	if handler == nil {
		handler = NopValuedHandler[T]{}
	}
	cooldown.handler.Store(handler)
}

// Active returns true if cooldown is currently active.
func (cooldown *Valued[T]) Active() bool {
	return cooldown.basic.Active()
}

// Remaining returns the duration until cooldown expiration.
func (cooldown *Valued[T]) Remaining() time.Duration {
	return cooldown.basic.Remaining()
}

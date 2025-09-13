package cooldown

import (
	"sync/atomic"
	"time"
)

// Valued represents basic cooldown with renew ability and custom values.
//
// Point of values is to use it in transactions and allow other users to handle
// the context value easily. All actions of the cooldown (like stop, start,
// renew) accepts a value. My personal opinion is to use in dragonfly
// (github.com/df-mc/dragonfly) world transactions. If you don't want to use
// values, there is zero implementation: CoolDown.
type Valued[T any] struct {
	basic *Basic

	duration atomic.Pointer[time.Duration]
	handler  atomic.Pointer[ValuedHandler[T]]

	timer atomic.Pointer[time.Timer]
}

// NewValued creates new Valued cooldown. User can omit handler argument nil.
func NewValued[T any](handler ValuedHandler[T]) *Valued[T] {
	cd := &Valued[T]{basic: new(Basic)}
	if handler != nil {
		// Only store if not nil, otherwise it can panic
		cd.handler.Store(&handler)
	}
	return cd
}

// Renew renews the cooldown, if it is currently active.
func (cooldown *Valued[T]) Renew(val T) {
	if !cooldown.Active() {
		// Not active, cannot renew
		return
	}

	durationPtr := cooldown.duration.Load()
	if durationPtr == nil {
		// Failed to load duration
		return
	}
	dur := *durationPtr
	if dur <= 0 {
		return
	}

	ctx := createContext(cooldown)
	if cooldown.Handler().HandleRenew(ctx, val); ctx.Cancelled() {
		return
	}

	cooldown.basic.Set(dur)
	timer := cooldown.timer.Load()
	if timer == nil {
		panic("tried to renew while timer is nil")
	}

	timer.Reset(dur)
}

// Start starts the cooldown, if it is not currently active.
func (cooldown *Valued[T]) Start(dur time.Duration, val T) {
	if cooldown.Active() {
		// Already active, cannot start again
		return
	}

	ctx := createContext(cooldown)
	if cooldown.Handler().HandleStart(ctx, val); ctx.Cancelled() {
		return
	}

	cooldown.duration.Store(&dur)
	cooldown.basic.Set(dur)

	if cooldown.timer.Load() != nil {
		panic("tried to start while timer is not nil")
	}

	cooldown.timer.Store(time.AfterFunc(dur, func() {
		cooldown.stop(ErrStopCauseExpired, true)
	}))
}

// Stop stops the cooldown, if it is currently active.
func (cooldown *Valued[T]) Stop(val T) {
	if !cooldown.Active() {
		// Not active, cannot stop
		return
	}

	cause := ErrStopCauseCancelled
	cooldown.Handler().HandleStop(cooldown, cause, val)
	cooldown.stop(cause, false) // already handled
}

func (cooldown *Valued[T]) stop(cause StopCause, handle bool) {
	if handle {
		var zeroT T
		cooldown.Handler().HandleStop(cooldown, cause, zeroT)
	}

	var dur time.Duration = -1
	cooldown.duration.Store(&dur)
	cooldown.basic.Reset()

	// Resetting timer
	timer := cooldown.timer.Load()
	if timer == nil {
		panic("tried to stop while timer is nil")
	}

	timer.Stop()
	cooldown.timer.Store(nil)
}

// Handler returns current cooldown handler. If it is not set, NopHandler will
// be returned.
func (cooldown *Valued[T]) Handler() ValuedHandler[T] {
	h := cooldown.handler.Load()
	if h == nil || *h == nil {
		var nop ValuedHandler[T] = NopValuedHandler[T]{}
		h = &nop
	}
	return *h
}

// Handle updates current cooldown handler. If user entered nil as argument,
// current handler will be removed.
func (cooldown *Valued[T]) Handle(handler ValuedHandler[T]) {
	// Making sure handler is never nil
	if handler == nil {
		// If it is nil, just updating it to no-operation handler
		handler = NopValuedHandler[T]{}
	}
	cooldown.handler.Store(&handler)
}

// Duration returns duration of the cooldown.
func (cooldown *Valued[T]) Duration() time.Duration {
	var dur time.Duration
	if v := cooldown.duration.Load(); v != nil {
		val := *v
		dur = val
	}
	return dur
}

// Active returns true if cooldown is currently active.
func (cooldown *Valued[T]) Active() bool {
	return cooldown.basic.Active()
}

// Remaining returns the duration until cooldown expiration.
func (cooldown *Valued[T]) Remaining() time.Duration {
	return cooldown.basic.Remaining()
}

// TODO context support

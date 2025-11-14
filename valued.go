package cooldown

import (
	"log/slog"
	"sync/atomic"
	"time"
)

var (
	// SafeMode is a global atomic boolean, marks if there shouldn't be any
	// potential panics (for prod).
	// If false, it can panic if something is very goes wrong.
	SafeMode = func() (b atomic.Bool) {
		b.Store(true)
		return
	}()
	// Logger is global instance of slog logger. It is used to log any warns
	// or errors.
	// It can be changed by user.
	Logger = func() (p atomic.Pointer[slog.Logger]) {
		p.Store(slog.Default())
		return
	}()
)

// Valued represents basic cooldown with renew ability and custom values.
// Point of values is to use them as transaction or context value.
type Valued[T any] struct {
	// basic is the underlying basic cooldown. It is used to control cooldown
	// basically (Set, Reset, Remaining).
	basic *Basic
	// duration is pointer to current duration of the cooldown.
	// It is used to Renew calls.
	duration atomic.Pointer[time.Duration]
	// handler is pointer to cooldown handler.
	// If nil, NopHandler will be used.
	handler atomic.Pointer[ValuedHandler[T]]
	// timer is pointer to active cooldown timer. This timer is created by
	// time.AfterFunc and is used to stop the cooldown.
	// If timer will expire (time.AfterFunc), it'll call stop function by
	// itself.
	timer atomic.Pointer[time.Timer]
}

// NewValued creates new Valued cooldown. User can omit handler argument nil.
func NewValued[T any](opts ...ValuedOption[T]) *Valued[T] {
	cd := &Valued[T]{basic: new(Basic)}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cd)
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
		if !SafeMode.Load() {
			panic(ErrRenewTimerNil)
		}
		Logger.Load().Warn("an error occurred while renewing cooldown", "err", ErrRenewTimerNil)
		// Starting new timer, that expires the cooldown
		cooldown.timer.Store(time.AfterFunc(dur, cooldown.expire))
		return
	}

	// Reset (renew) the timer with new duration.
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

	if t := cooldown.timer.Load(); t != nil {
		if !SafeMode.Load() {
			panic(ErrStartTimerNotNil)
		}
		// Log the error on safe mode
		Logger.Load().Warn("an error occurred while starting cooldown", "err", ErrStartTimerNotNil)
		// And clear actual timer
		t.Stop()
		cooldown.timer.Store(nil)
	}

	// Starting timer, that expires the cooldown
	cooldown.timer.Store(time.AfterFunc(dur, cooldown.expire))
}

func (cooldown *Valued[T]) expire() {
	var zeroT T
	cooldown.Handler().HandleStop(cooldown, ErrStopCauseExpired, zeroT)
	cooldown.stop()
}

// Stop stops the cooldown, if it is currently active.
func (cooldown *Valued[T]) Stop(val T) {
	if !cooldown.Active() {
		// Not active, cannot stop
		return
	}

	cause := ErrStopCauseCancelled
	cooldown.Handler().HandleStop(cooldown, cause, val)
	cooldown.stop()
}

func (cooldown *Valued[T]) stop() {
	cooldown.duration.Store(nil)
	cooldown.basic.Reset()

	timer := cooldown.timer.Load()
	if timer == nil {
		if !SafeMode.Load() {
			panic(ErrStopTimerNil)
		}
		Logger.Load().Warn("an error occurred while stopping cooldown", "err", ErrStopTimerNil)
		return
	}

	// Stop the timer to prevent cooldown expiration
	timer.Stop()
	// Invalidate the timer
	cooldown.timer.Store(nil)
}

// Handler returns current cooldown handler. If it is not set, NopHandler will
// be returned.
func (cooldown *Valued[T]) Handler() ValuedHandler[T] {
	handlerPtr := cooldown.handler.Load()
	if handlerPtr == nil {
		// Nil pointer (no handler set)
		return NopValuedHandler[T]{}
	}
	handlerVal := *handlerPtr
	if handlerVal == nil {
		// Nil handler, but not nil pointer ???
		cooldown.handler.Store(nil) // Store nil pointer to optimize (get nop handler instantly)
		return NopValuedHandler[T]{}
	}
	return handlerVal
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
func (cooldown *Valued[T]) Duration() (dur time.Duration) {
	if v := cooldown.duration.Load(); v != nil {
		// Can safely dereference v, even if duration is not set, it is not
		// pointer, so it'll automatically be 0
		val := *v
		dur = val
	}
	return
}

// Active returns true if cooldown is currently active.
func (cooldown *Valued[T]) Active() bool {
	return cooldown.basic.Active()
}

// Remaining returns the duration until cooldown expiration.
// If cooldown is not active, it returns number that either zero or negative.
func (cooldown *Valued[T]) Remaining() time.Duration {
	return cooldown.basic.Remaining()
}

// Basic returns underlying basic cooldown of this valued cooldown.
// WARNING: It should be only used for testing purposes.
func (cooldown *Valued[T]) Basic() *Basic {
	return cooldown.basic
}

// TODO context support

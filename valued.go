package cooldown

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/k4ties/cooldown/internal/event"
)

// Valued represents cooldown with renew ability and values.
// Point of values is to use them as transaction or context value.
type Valued[T any] struct {
	mu sync.RWMutex // also controls basic

	basic    *Basic
	duration time.Duration
	timer    *time.Timer

	handler atomic.Pointer[ValuedHandler[T]]
}

// NewValued creates new Valued cooldown.
func NewValued[T any](opts ...ValuedOption[T]) *Valued[T] {
	cd := &Valued[T]{basic: new(Basic)}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cd)
	}
	if cd.handler.Load() == nil {
		h := ValuedHandler[T](NopValuedHandler[T]{})
		cd.handler.Store(&h)
	}
	return cd
}

// Renew ...
func (cooldown *Valued[T]) Renew(val T) {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()
	cooldown.RenewUnsafe(val)
}

func (cooldown *Valued[T]) RenewUnsafe(val T) {
	if !cooldown.ActiveUnsafe() {
		return
	}
	dur := cooldown.duration
	if dur <= 0 {
		return
	}
	timer := cooldown.timer
	if timer == nil {
		return
	}
	ctx := event.C(cooldown)
	if cooldown.Handler().HandleRenew(ctx, dur, val); ctx.Cancelled() {
		return
	}
	cooldown.basic.SetUnsafe(dur)
	timer.Reset(dur)
}

// Start ...
func (cooldown *Valued[T]) Start(dur time.Duration, val T) bool {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()
	return cooldown.StartUnsafe(dur, val)
}

func (cooldown *Valued[T]) StartUnsafe(dur time.Duration, val T) bool {
	if dur <= 0 {
		return false
	}
	if cooldown.ActiveUnsafe() {
		cooldown.StopUnsafe(val)
	}
	ctx := event.C(cooldown)
	if cooldown.Handler().HandleStart(ctx, dur, val); ctx.Cancelled() {
		return false
	}
	cooldown.duration = dur
	cooldown.timer = time.AfterFunc(dur, cooldown.expire)
	cooldown.basic.SetUnsafe(dur)
	return true
}

func (cooldown *Valued[T]) expire() {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()

	var zeroT T
	cooldown.Handler().HandleStop(cooldown, ErrStopCauseExpired, zeroT)
	cooldown.doStopUnsafe(zeroT)
}

// Stop ...
func (cooldown *Valued[T]) Stop(val T) {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()
	cooldown.StopUnsafe(val)
}

func (cooldown *Valued[T]) StopUnsafe(val T) {
	if !cooldown.ActiveUnsafe() {
		return
	}
	cooldown.Handler().HandleStop(cooldown, ErrStopCauseCancelled, val)
	cooldown.doStopUnsafe(val)
}

func (cooldown *Valued[T]) doStopUnsafe(val T) {
	if cooldown.PausedUnsafe() {
		cooldown.doResumeUnsafe(val, false) // we will stop the timer
	}
	cooldown.duration = 0
	cooldown.basic.ResetUnsafe()

	if timer := cooldown.timer; timer != nil {
		timer.Stop()
		cooldown.timer = nil
	}
}

// Pause ...
func (cooldown *Valued[T]) Pause(val T) bool {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()
	return cooldown.PauseUnsafe(val)
}

func (cooldown *Valued[T]) PauseUnsafe(val T) bool {
	if cooldown.PausedUnsafe() || !cooldown.ActiveUnsafe() {
		return false
	}
	timer := cooldown.timer
	if timer == nil {
		return false
	}
	ctx := event.C(cooldown)
	if cooldown.Handler().HandlePause(ctx, val); ctx.Cancelled() {
		return false
	}
	if !cooldown.basic.PauseUnsafe() {
		return false
	}
	ok := timer.Stop()
	cooldown.timer = nil // Resume will create new timer
	return ok
}

// Resume ...
func (cooldown *Valued[T]) Resume(val T) bool {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()
	return cooldown.ResumeUnsafe(val)
}

func (cooldown *Valued[T]) ResumeUnsafe(val T) bool {
	return cooldown.doResumeUnsafe(val, true)
}

func (cooldown *Valued[T]) doResumeUnsafe(val T, resetTimer bool) bool {
	if !cooldown.PausedUnsafe() {
		return false
	}
	dur := cooldown.duration
	if dur <= 0 {
		return false
	}
	ctx := event.C(cooldown)
	if cooldown.Handler().HandleResume(ctx, val); ctx.Cancelled() {
		return false
	}
	if !cooldown.basic.ResumeUnsafe() {
		return false
	}
	if resetTimer {
		// RemainingUnsafe also accounts for paused state
		cooldown.timer = time.AfterFunc(cooldown.RemainingUnsafe(), cooldown.expire)
	}
	return true
}

// TogglePause ...
func (cooldown *Valued[T]) TogglePause(val T) bool {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()
	return cooldown.TogglePauseUnsafe(val)
}

func (cooldown *Valued[T]) TogglePauseUnsafe(val T) bool {
	if cooldown.PausedUnsafe() {
		return cooldown.ResumeUnsafe(val)
	}
	return cooldown.PauseUnsafe(val)
}

// Handler ...
func (cooldown *Valued[T]) Handler() ValuedHandler[T] {
	// if properly initialized this is never nil
	return *cooldown.handler.Load()
}

// Handle ...
func (cooldown *Valued[T]) Handle(handler ValuedHandler[T]) {
	if handler == nil {
		handler = NopValuedHandler[T]{}
	}
	cooldown.handler.Store(&handler)
}

// Duration ...
func (cooldown *Valued[T]) Duration() time.Duration {
	cooldown.mu.RLock()
	defer cooldown.mu.RUnlock()
	return cooldown.DurationUnsafe()
}

func (cooldown *Valued[T]) DurationUnsafe() time.Duration {
	return cooldown.duration
}

// Active ...
func (cooldown *Valued[T]) Active() bool {
	cooldown.mu.RLock()
	defer cooldown.mu.RUnlock()
	return cooldown.ActiveUnsafe()
}

func (cooldown *Valued[T]) ActiveUnsafe() bool {
	return cooldown.basic.ActiveUnsafe()
}

// Remaining ...
func (cooldown *Valued[T]) Remaining() time.Duration {
	cooldown.mu.RLock()
	defer cooldown.mu.RUnlock()
	return cooldown.RemainingUnsafe()
}

func (cooldown *Valued[T]) RemainingUnsafe() time.Duration {
	return cooldown.basic.RemainingUnsafe()
}

// Paused ...
func (cooldown *Valued[T]) Paused() bool {
	cooldown.mu.RLock()
	defer cooldown.mu.RUnlock()
	return cooldown.PausedUnsafe()
}

func (cooldown *Valued[T]) PausedUnsafe() bool {
	return cooldown.basic.PausedUnsafe()
}

func (cooldown *Valued[T]) L() *sync.RWMutex {
	return &cooldown.mu
}

// TODO context.Context support ? context.AfterFunc

// TODO return errors , maybe create something like 'MustPause()' methods with logger

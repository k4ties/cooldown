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
func (cooldown *Valued[T]) Start(dur time.Duration, val T) {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()
	cooldown.StartUnsafe(dur, val)
}

func (cooldown *Valued[T]) StartUnsafe(dur time.Duration, val T) {
	if cooldown.ActiveUnsafe() || dur <= 0 {
		return
	}
	ctx := event.C(cooldown)
	if cooldown.Handler().HandleStart(ctx, dur, val); ctx.Cancelled() {
		return
	}
	cooldown.duration = dur
	cooldown.timer = time.AfterFunc(dur, cooldown.expire)
	cooldown.basic.SetUnsafe(dur)
}

func (cooldown *Valued[T]) expire() {
	cooldown.mu.Lock()
	defer cooldown.mu.Unlock()

	var zeroT T
	cooldown.Handler().HandleStop(cooldown, ErrStopCauseExpired, zeroT)
	cooldown.doStopUnsafe()
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
	cooldown.doStopUnsafe()
}

func (cooldown *Valued[T]) doStopUnsafe() {
	cooldown.duration = 0
	cooldown.basic.ResetUnsafe()

	if timer := cooldown.timer; timer != nil {
		timer.Stop()
		cooldown.timer = nil
	}
}

// Handler ...
func (cooldown *Valued[T]) Handler() ValuedHandler[T] {
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

// TODO context.Context support ?

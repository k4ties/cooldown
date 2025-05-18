package cooldown

import (
	"context"
	"github.com/k4ties/cooldown/internal/event"
	"sync"
	"sync/atomic"
	"time"
)

// WithVal represents a cooldown with per-tick handler and renew function.
type WithVal[T any] struct {
	exp   atomic.Pointer[time.Time]
	renew atomic.Pointer[chan struct{}]

	cancel atomic.Pointer[context.CancelCauseFunc]

	handler atomic.Pointer[HandlerWithVal[T]]
	wg      sync.WaitGroup

	taskFunc     StartTaskFunc
	tickerActive atomic.Bool
}

// NewWithVal returns new blank cooldown with value.
func NewWithVal[T any](h HandlerWithVal[T], opts ...Option[T]) *WithVal[T] {
	cd := &WithVal[T]{}
	cd.Handle(h)
	cd.cancel.Store(nil)
	cd.taskFunc = cd.startTick
	for _, opt := range opts {
		opt(cd)
	}
	return cd
}

func (c *WithVal[T]) Start(dur time.Duration, val T) {
	if c.Active() {
		c.Stop(val)
	}

	ctx := event.Ctx(c)
	if c.Handler().HandleStart(convContext(ctx), val); ctx.Cancelled() {
		return
	}

	c.set(dur)
}

// set sets the cooldown duration to the specified one. If cooldown is active, it will be stopped.
func (c *WithVal[T]) set(dur time.Duration) {
	c.startTickTask(dur)

	exp := time.Now().Add(dur)
	c.exp.Store(&exp)
}

// Renew renews the CoolDown. If it's not currently active, it'll panic.
func (c *WithVal[T]) Renew(val T) {
	if !c.hasRenewChan() || !c.Active() {
		panic("unable to renew")
	}

	ctx := event.Ctx(c)
	if c.Handler().HandleRenew(convContext(ctx), val); ctx.Cancelled() {
		return
	}

	c.renewChanWrite() <- struct{}{}
}

// Remaining returns time until expiration of the CoolDown.
func (c *WithVal[T]) Remaining() time.Duration {
	exp := c.exp.Load()
	if exp == nil || (*exp).Equal(time.Time{}) {
		return -1
	}
	return time.Until(*exp)
}

// Stop resets current cooldown. If currently CoolDown ticker is active, it will be
// stopped immediately.
func (c *WithVal[T]) Stop(val T) {
	if !c.Active() {
		panic("trying to reset while cooldown is not active")
	}
	c.reset(StopCauseCancelled, val)
}

func (c *WithVal[T]) reset(cause StopCause, val any) {
	if cancel := c.getCancel(); cancel != nil {
		cancel(cause)

		var userVal T
		if val != nil {
			userVal, _ = val.(T)
		}

		c.Handler().HandleStop(c, cause, userVal)
	}

	if renewChan := c.renewChanWrite(); c.hasRenewChan() {
		close(renewChan)
		c.renew.Store(nil)
	}

	c.exp.Store(nil)
	c.wg.Wait()
}

// Active returns true if cooldown is currently active.
func (c *WithVal[T]) Active() bool {
	exp := c.exp.Load()
	if exp == nil || exp.Equal(time.Time{}) {
		return false
	}
	return exp.After(time.Now())
}

// Handle sets new handler to the cooldown. If entered handler is nil, current handler
// will be removed.
func (c *WithVal[T]) Handle(h HandlerWithVal[T]) {
	if h == nil {
		h = NopHandler[T]{}
	}
	c.handler.Store(&h)
}

// Handler returns current cooldown handler.
func (c *WithVal[T]) Handler() HandlerWithVal[T] {
	return *c.handler.Load()
}

// getCancel ...
func (c *WithVal[T]) getCancel() context.CancelCauseFunc {
	cancel := c.cancel.Load()
	if cancel == nil {
		return nil
	}
	return *cancel
}

// renewChan ...
func (c *WithVal[T]) renewChan() chan struct{} {
	val := c.renew.Load()
	if val == nil {
		return nil
	}
	return *val
}

// renewChanWrite ...
func (c *WithVal[T]) renewChanWrite() chan<- struct{} {
	return c.renewChan()
}

// renewChanWrite ...
func (c *WithVal[T]) renewChanRead() <-chan struct{} {
	return c.renewChan()
}

// hasRenewChan ...
func (c *WithVal[T]) hasRenewChan() bool {
	return c.renewChan() != nil
}

// SetCancel ...
func (c *WithVal[T]) SetCancel(val context.CancelCauseFunc) {
	c.cancel.Store(&val)
}

// SetTickerActive ...
func (c *WithVal[T]) SetTickerActive(status bool) {
	c.tickerActive.Store(status)
}

// convContext ...
func convContext[T any](ctx *event.Context[*WithVal[T]]) *ContextWithVal[T] {
	return (*ContextWithVal[T])(ctx)
}

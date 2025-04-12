package cooldown

import (
	"context"
	"github.com/df-mc/dragonfly/server/event"
	"sync/atomic"
	"time"
)

// CoolDown represents a cooldown with per-tick handler.
type CoolDown struct {
	exp atomic.Pointer[time.Time]

	cancel  atomic.Pointer[context.CancelCauseFunc]
	handler atomic.Pointer[Handler]

	renew atomic.Pointer[chan struct{}]
}

// New returns new blank cooldown.
func New(h Handler) *CoolDown {
	if h == nil {
		h = NopHandler{}
	}

	cd := &CoolDown{}
	cd.Handle(h)
	cd.cancel.Store(&zeroCancel)

	return cd
}

// zeroCancel ...
var zeroCancel context.CancelCauseFunc

// Set sets the cooldown duration to the specified one. If cooldown is active, it will be stopped.
func (c *CoolDown) Set(dur time.Duration) {
	ctx := event.C(c)
	if c.Handler().HandleStart(ctx); ctx.Cancelled() {
		return
	}
	if c.Active() {
		c.Reset()
	}

	c.startTick(dur, context.Background())

	exp := time.Now().Add(dur)
	c.exp.Store(&exp)
}

// Renew renews the CoolDown. If it's not currently active, it'll panic.
func (c *CoolDown) Renew() {
	if !c.hasRenewChan() || !c.Active() {
		panic("unable to renew")
	}

	ctx := event.C(c)
	if c.Handler().HandleRenew(ctx); ctx.Cancelled() {
		return
	}

	c.renewChan() <- struct{}{}
}

// Remaining returns time until expiration of the CoolDown.
func (c *CoolDown) Remaining() time.Duration {
	exp := c.exp.Load()
	return time.Until(*exp)
}

// Reset resets current cooldown. If currently CoolDown ticker is active, it will be
// stopped immediately.
func (c *CoolDown) Reset() {
	if !c.Active() {
		panic("trying to reset while cooldown is not active")
	}
	c.reset(StopCauseCancelled)
}

func (c *CoolDown) reset(cause StopCause) {
	if cancel := c.getCancel(); cancel != nil {
		cancel(cause)
		c.cancel.Store(&zeroCancel)
	}
	if renewChan := c.renewChan(); renewChan != nil {
		close(renewChan)
		c.renew.Store(new(chan struct{}))
	}
	c.exp.Store(&time.Time{})
}

// Active returns true if cooldown is currently active.
func (c *CoolDown) Active() bool {
	exp := c.exp.Load()
	if exp == nil {
		return false
	}
	return (*exp).After(time.Now())
}

// Handle sets new handler to the cooldown. If this handler is nil, handler will
// become no-operation.
func (c *CoolDown) Handle(h Handler) {
	if h == nil {
		h = NopHandler{}
	}
	c.handler.Store(&h)
}

// Handler returns current cooldown handler.
func (c *CoolDown) Handler() Handler {
	val := c.handler.Load()
	return (*val).(Handler)
}

// getCancel ...
func (c *CoolDown) getCancel() context.CancelCauseFunc {
	return *c.cancel.Load()
}

// hasRenewChan ...
func (c *CoolDown) hasRenewChan() bool {
	return c.renewChan() != nil
}

// renewChan ...
func (c *CoolDown) renewChan() chan<- struct{} {
	val := c.renew.Load()
	if val == nil {
		return nil
	}
	return *val
}

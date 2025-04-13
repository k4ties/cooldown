package cooldown

import (
	"context"
	"github.com/df-mc/dragonfly/server/event"
	"sync"
	"sync/atomic"
	"time"
)

// CoolDown represents a cooldown with per-tick handler and renew function.
type CoolDown struct {
	exp atomic.Value //time.Time

	cancel  atomic.Value //context.CancelCauseFunc
	handler atomic.Value //Handler

	wg sync.WaitGroup

	renew atomic.Value //chan struct{}
}

// New returns new blank cooldown.
func New(h Handler) *CoolDown {
	if h == nil {
		h = NopHandler{}
	}

	cd := &CoolDown{}
	cd.Handle(h)
	cd.cancel.Store(zeroCancel)

	return cd
}

// zeroCancel ...
var zeroCancel context.CancelCauseFunc

func (c *CoolDown) Start(dur time.Duration) {
	if c.Active() {
		c.Stop()
	}
	c.set(dur)
}

// set sets the cooldown duration to the specified one. If cooldown is active, it will be stopped.
func (c *CoolDown) set(dur time.Duration) {
	ctx := event.C(c)
	if c.Handler().HandleStart(ctx); ctx.Cancelled() {
		return
	}

	c.startTick(dur, context.Background())
	c.exp.Store(time.Now().Add(dur))
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
	if exp == nil {
		return -1
	}
	return time.Until(exp.(time.Time))
}

// Stop resets current cooldown. If currently CoolDown ticker is active, it will be
// stopped immediately.
func (c *CoolDown) Stop() {
	if !c.Active() {
		panic("trying to reset while cooldown is not active")
	}
	c.reset(StopCauseCancelled)
}

func (c *CoolDown) reset(cause StopCause) {
	if cancel := c.getCancel(); cancel != nil {
		cancel(cause)
		c.cancel.Store(zeroCancel)
	}
	if renewChan := c.renewChan(); renewChan != nil {
		close(renewChan)
		c.renew.Store(*new(chan struct{}))
	}
	c.exp.Store(time.Time{})
	c.wg.Wait()
}

// Active returns true if cooldown is currently active.
func (c *CoolDown) Active() bool {
	exp := c.exp.Load()
	if exp == nil {
		return false
	}
	return exp.(time.Time).After(time.Now())
}

// Handle sets new handler to the cooldown. If this handler is nil, handler will
// become no-operation.
func (c *CoolDown) Handle(h Handler) {
	if h == nil {
		h = NopHandler{}
	}
	c.handler.Store(h)
}

// Handler returns current cooldown handler.
func (c *CoolDown) Handler() Handler {
	val := c.handler.Load()
	if val == nil {
		return nil
	}
	return val.(Handler)
}

// getCancel ...
func (c *CoolDown) getCancel() context.CancelCauseFunc {
	val := c.cancel.Load()
	if val == nil {
		return nil
	}
	return val.(context.CancelCauseFunc)
}

// renewChan ...
func (c *CoolDown) renewChan() chan<- struct{} {
	val := c.renew.Load()
	if val == nil {
		return nil
	}
	return val.(chan struct{})
}

// hasRenewChan ...
func (c *CoolDown) hasRenewChan() bool {
	return c.renewChan() != nil
}

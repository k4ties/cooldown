package cooldown

import (
	"context"
	"github.com/df-mc/dragonfly/server/event"
	"go.uber.org/atomic"
	"time"
)

// CoolDown represents a cooldown with per-tick handler.
type CoolDown struct {
	exp atomic.Time

	cancel,
	handler atomic.Value
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
var zeroCancel context.CancelFunc

// Set sets the cooldown duration to the specified one. If cooldown is active, it will be stopped.
func (c *CoolDown) Set(dur time.Duration) {
	if !c.reHandleSet() {
		return
	}
	if c.Active() {
		c.Reset()
	}

	c.startTick(dur, context.Background())

	exp := time.Now().Add(dur)
	c.exp.Store(exp)
}

// reHandleSet handles renew if cooldown is currently active, otherwise it will handle start. It
// returns false if context is cancelled and true if it's not.
func (c *CoolDown) reHandleSet() bool {
	ctx := event.C(c)
	handler := c.Handler()

	if c.Active() {
		// if currently active, it is renewed event
		if handler.HandleRenew(ctx); ctx.Cancelled() {
			return false
		}
	} else {
		// otherwise it is start event
		if handler.HandleStart(ctx); ctx.Cancelled() {
			return false
		}
	}
	return true
}

// Remaining returns time until expiration of the CoolDown.
func (c *CoolDown) Remaining() time.Duration {
	exp := c.exp.Load()
	return time.Until(exp)
}

// Reset resets current cooldown. If currently CoolDown ticker is active, it will be
// stopped immediately.
func (c *CoolDown) Reset() {
	if cancel := c.getCancel(); cancel != nil {
		cancel()
		c.cancel.Store((context.CancelFunc)(nil))
	}
	c.exp.Store(time.Time{})
}

// Active returns true if cooldown is currently active.
func (c *CoolDown) Active() bool {
	return c.exp.Load().After(time.Now())
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
	return val.(Handler)
}

// getCancel ...
func (c *CoolDown) getCancel() context.CancelFunc {
	return c.cancel.Load().(context.CancelFunc)
}

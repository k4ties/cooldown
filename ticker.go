package cooldown

import (
	"context"
	"time"
)

// TicksPerSecond is the N amount that is used to calculate tick duration. It is
// used in the CoolDown ticker.
var TicksPerSecond time.Duration = 20

// tickDuration ...
func tickDuration() time.Duration {
	return time.Second / TicksPerSecond
}

// startTick starts the ticker task of the CoolDown.
func (c *CoolDown) startTick(dur time.Duration, parent context.Context) {
	ctx, cancel := context.WithCancelCause(parent)
	c.cancel.Store(cancel)

	ticker := time.NewTicker(tickDuration())
	timer := time.NewTimer(dur)

	go c.tickTask(ctx, ticker, timer, dur)
}

// tickTask start the main ticker task of the CoolDown.
func (c *CoolDown) tickTask(ctx context.Context, ticker *time.Ticker, timer *time.Timer, dur time.Duration) {
	c.wg.Add(1)
	defer c.wg.Done()

	defer func() {
		ticker.Stop()
		timer.Stop()
	}()

	var cause StopCause
	var tick int64

	defer func() {
		c.cancel.Store(zeroCancel)
		c.Handler().HandleStop(c, cause)
	}()

	renew := make(chan struct{})

	if c.hasRenewChan() {
		panic("tried to start ticking when already ticking")
	}

	c.renew.Store(renew)

	for {
		select {
		case <-ticker.C:
			tick++
			c.Handler().HandleTick(c, tick)
		case <-timer.C:
			cause = StopCauseExpired
			return
		case <-(<-chan struct{})(renew): // only receive
			ticker.Reset(tickDuration())
			timer.Reset(dur)

			c.exp.Store(time.Now().Add(dur))
		case <-ctx.Done():
			cause = StopCauseCancelled
			return
		}
	}
}

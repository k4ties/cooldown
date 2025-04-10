package cooldown

import (
	"context"
	"errors"
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
	c.cancel.Store(&cancel)

	ticker := time.NewTicker(tickDuration())
	timer := time.NewTimer(dur)

	go c.tick(ctx, ticker, timer)
}

// tick start the main ticker task of the CoolDown.
func (c *CoolDown) tick(ctx context.Context, ticker *time.Ticker, timer *time.Timer) {
	defer func() {
		ticker.Stop()
		timer.Stop()
	}()

	var cause StopCause
	var tick int64

	defer func() {
		c.cancel.Store(&zeroCancel)

		if err := context.Cause(ctx); err != nil {
			// if it's renew, it is not stop event, so we'll return from func to prevent handling stop
			if errors.Is(err, StopCauseRenew) {
				return
			}
		}

		c.Handler().HandleStop(c, cause)
	}()

	for {
		select {
		case <-ticker.C:
			tick++
			c.Handler().HandleTick(c, tick)
		case <-timer.C:
			cause = StopCauseExpired
			return
		case <-ctx.Done():
			cause = StopCauseCancelled
			return
		}
	}
}

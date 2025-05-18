package cooldown

import (
	"context"
	"sync"
	"time"
)

// DefaultTicksPerSecond is the N amount that is used to calculate tick duration. It is
// used in the CoolDown ticker.
var DefaultTicksPerSecond = 20

// tickDuration ...
func tickDuration() time.Duration {
	return time.Second / time.Duration(DefaultTicksPerSecond)
}

type StartTaskFunc = func(*TickData)

func (c *WithVal[T]) startTickTask(dur time.Duration) {
	if c.hasRenewChan() {
		panic("tried to start ticking when renew chan isn't removed")
	}

	renew := make(chan struct{})
	c.renew.Store(&renew)

	timer := time.NewTimer(dur)
	var tick int64

	data := TickData{
		WaitGroup: &c.wg,
		Timer:     timer,
		Duration:  dur,
		TickPtr:   &tick,
	}

	c.wg.Add(1)
	c.taskFunc(&data)
}

// startTick starts the ticker task of the CoolDown.
func (c *WithVal[T]) startTick(data *TickData) {
	ctx, cancel := context.WithCancelCause(context.Background())
	data.Context = ctx
	c.cancel.Store(&cancel)

	go func() {
		c.SetTickerActive(true)

		ticker := time.NewTicker(tickDuration())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.Tick(data, *new(T))
			}
		}
	}()
}

type TickData struct {
	Context   context.Context
	WaitGroup *sync.WaitGroup
	Timer     *time.Timer
	Duration  time.Duration
	TickPtr   *int64
}

// Tick ...
func (c *WithVal[T]) Tick(data *TickData, val T) {
	if c.tickerActive.Load() {
		return
	}

	stopRoutine := func() {
		c.SetTickerActive(false)
		data.Timer.Stop()
		data.WaitGroup.Done()
	}

	select {
	case <-data.Timer.C:
		c.reset(StopCauseExpired, nil)
		stopRoutine()
	case <-data.Context.Done():
		stopRoutine()
	case <-c.renewChanRead():
		data.Timer.Reset(data.Duration)

		exp := time.Now().Add(data.Duration)
		c.exp.Store(&exp)
	default:
		*data.TickPtr++
		c.Handler().HandleTick(c, *data.TickPtr, val)
	}
}

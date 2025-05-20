package cooldown

import (
	"context"
	"github.com/k4ties/cooldown/internal/atomic"
	"github.com/k4ties/gq"
	"io"
	"sync"
	"time"
)

// Proc is current cooldown Processor.
var Proc = NewDefaultProcessor()

func init() {
	ctx := context.Background()

	proc := Proc
	go proc.StartTracking(ctx, tickDuration())
}

// Processor is used to track processable(s) expiration(s).
type Processor interface {
	io.Closer
	// StartTracking start tracking cooldown expirations.
	StartTracking(parent context.Context, every time.Duration)
	// Append appends processable object to the Processor.
	Append(p processable)
	// Remove removes processable object from the Processor, if it existed.
	Remove(p processable)
	// Running returns true, if processor is currently running.
	Running() bool
}

// processor is Processor implementation.
type processor struct {
	running atomic.Value[bool]
	close   sync.Once

	cooldowns   gq.Set[processable]
	cooldownsMu sync.RWMutex

	cancel atomic.Value[context.CancelFunc]
	ticker atomic.Value[*time.Ticker]
}

// NewDefaultProcessor creates new default Processor impl.
func NewDefaultProcessor() Processor {
	proc := &processor{cooldowns: make(gq.Set[processable])}
	proc.cancel = atomic.NewValue[context.CancelFunc]()
	proc.ticker = atomic.NewValue[*time.Ticker]()
	proc.running = atomic.NewValue[bool]()
	return proc
}

// StartTracking starts tracking of processable(s) expiration.
func (processor *processor) StartTracking(parent context.Context, every time.Duration) {
	ctx, cancel := context.WithCancel(parent)

	processor.cancel.Store(cancel)
	processor.ticker.Store(time.NewTicker(every))

	processor.running.Store(true)
	defer processor.running.Store(false)

	for {
		ticker, ok := processor.ticker.Load()
		if !ok {
			return
		}

		select {
		case <-ticker.C:
			func() {
				processor.cooldownsMu.RLock()
				defer processor.cooldownsMu.RUnlock()

				for cooldown := range processor.cooldowns {
					// checking if cooldown is expired
					if expiration := getExpiration(cooldown); expiration.Before(time.Now()) {
						// expired
						cooldown.stop(StopCauseExpired, true)
					}
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}

// Append ...
func (processor *processor) Append(p processable) {
	processor.cooldownsMu.Lock()
	if !processor.cooldowns.Contains(p) {
		processor.cooldowns.Add(p)
	}
	processor.cooldownsMu.Unlock()
}

// Remove ...
func (processor *processor) Remove(p processable) {
	processor.cooldownsMu.Lock()
	if processor.cooldowns.Contains(p) {
		processor.cooldowns.Delete(p)
	}
	processor.cooldownsMu.Unlock()
}

// Running should return true, if Processor is currently running.
func (processor *processor) Running() bool {
	v, _ := processor.running.Load()
	return v == true
}

// Close ...
func (processor *processor) Close() error {
	processor.close.Do(func() {
		processor.cooldownsMu.Lock()
		defer processor.cooldownsMu.Unlock()

		// If it is not running, there's nothing to do, excluding clearing set.
		if !processor.Running() {
			processor.cooldowns.Clear()
			return
		}

		// Calling cancel function. If processor is running, it shouldn't be nil.
		cancel, _ := processor.cancel.Load()
		cancel()

		// Forcing all cooldowns to stop
		for cooldown := range processor.cooldowns {
			cooldown.stop(StopCauseClosed, true)
			// After cooldown is stopped, deleting it from the set
			processor.cooldowns.Delete(cooldown)
		}
	})
	return nil
}

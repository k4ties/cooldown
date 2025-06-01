package cooldown

import (
	"context"
	"github.com/k4ties/cooldown/internal/atomic"
	"github.com/k4ties/gq"
	"github.com/sasha-s/go-deadlock"
	"io"
	"sync"
	"time"
)

// Proc is current cooldown Processor.
var Proc = NewDefaultProcessor()

// Processor is used to track processable(s) expiration(s).
type Processor interface {
	io.Closer
	// StartTracking start tracking cooldown expirations.
	StartTracking(parent context.Context)
	// Append appends processable object to the Processor.
	Append(p Processable)
	// Remove removes processable object from the Processor, if it existed.
	Remove(p Processable)
	// Running returns true, if processor is currently running.
	Running() bool
}

// processor is Processor implementation.
type processor struct {
	running atomic.Value[bool]
	close   sync.Once

	cooldowns   gq.Set[Processable]
	cooldownsMu deadlock.RWMutex

	cancel atomic.Value[context.CancelFunc]
}

// NewDefaultProcessor creates new default Processor impl.
func NewDefaultProcessor() Processor {
	proc := &processor{cooldowns: make(gq.Set[Processable])}
	proc.cancel = atomic.NewValue[context.CancelFunc]()
	proc.running = atomic.NewValue[bool]()
	return proc
}

// StartTracking starts tracking of processable(s) expiration.
func (processor *processor) StartTracking(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	processor.cancel.Store(cancel)

	processor.running.Store(true)
	defer processor.running.Store(false)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			func() {
				processor.cooldownsMu.Lock()
				defer processor.cooldownsMu.Unlock()

				for cooldown := range processor.cooldowns {
					// checking if cooldown is expired
					if expiration := getExpiration(cooldown); expiration.Before(time.Now()) {
						// expired
						cooldown.UnsafeStop(StopCauseExpired, true)
						processor.remove(cooldown)
					}
				}
			}()
		}
	}
}

// Append ...
func (processor *processor) Append(p Processable) {
	processor.cooldownsMu.Lock()
	if !processor.cooldowns.Contains(p) {
		processor.cooldowns.Add(p)
	}
	processor.cooldownsMu.Unlock()
}

// Remove ...
func (processor *processor) Remove(p Processable) {
	processor.cooldownsMu.Lock()
	processor.remove(p)
	processor.cooldownsMu.Unlock()
}

// remove removes processable from the processor set, if it exists.
func (processor *processor) remove(p Processable) {
	if processor.cooldowns.Contains(p) {
		processor.cooldowns.Delete(p)
	}
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
			cooldown.UnsafeStop(StopCauseClosed, true)
			// After cooldown is stopped, deleting it from the set
			processor.cooldowns.Delete(cooldown)
		}
	})
	return nil
}

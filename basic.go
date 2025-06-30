package cooldown

import (
	"github.com/k4ties/cooldown/internal/atomic"
	"time"
)

// Basic represents very basic cooldown.
type Basic struct {
	expiration atomic.Value[time.Time]
}

// NewBasic creates new basic cooldown impl.
func NewBasic() *Basic {
	basic := &Basic{}
	basic.expiration = atomic.NewValue[time.Time]()
	return basic
}

// Set updates state of the cooldown.
func (cooldown *Basic) Set(dur time.Duration) {
	now := time.Now()
	expiration := now.Add(dur)
	cooldown.expiration.Store(expiration)
}

var zeroTime = time.Time{}

// Reset resets the cooldown expiration.
func (cooldown *Basic) Reset() {
	cooldown.expiration.Store(zeroTime)
}

// Active returns true if cooldown is currently active.
func (cooldown *Basic) Active() bool {
	expiration, ok := cooldown.expiration.Load()
	if !ok || expiration.Equal(zeroTime) {
		return false
	}
	now := time.Now()
	// If expiration date is before current date, it is expired. If it is not,
	// expiration date haven't been passed
	return !expiration.Before(now)
}

// Remaining returns the duration until cooldown expiration.
func (cooldown *Basic) Remaining() time.Duration {
	if !cooldown.Active() {
		return -1
	}
	expiration, _ := cooldown.expiration.Load()
	return time.Until(expiration)
}

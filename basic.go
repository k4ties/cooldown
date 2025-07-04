package cooldown

import (
	"sync/atomic"
	"time"
)

// Basic represents very basic cooldown.
type Basic struct {
	expiration atomic.Pointer[time.Time]
}

// Set updates state of the cooldown.
func (cooldown *Basic) Set(dur time.Duration) {
	now := time.Now()
	expiration := now.Add(dur)
	cooldown.expiration.Store(&expiration)
}

var zeroTime = time.Time{}

// Reset resets the cooldown expiration.
func (cooldown *Basic) Reset() {
	cooldown.expiration.Store(nil)
}

// Active returns true if cooldown is currently active.
func (cooldown *Basic) Active() bool {
	expiration := cooldown.expiration.Load()
	if expiration == nil || expiration.Equal(zeroTime) {
		return false
	}
	now := time.Now()
	// If expiration date is before current date, it is expired. If it is not,
	// expiration date haven't been passed
	return !(*expiration).Before(now)
}

// Remaining returns the duration until cooldown expiration.
func (cooldown *Basic) Remaining() time.Duration {
	if !cooldown.Active() {
		return -1
	}
	expiration := cooldown.expiration.Load()
	if expiration == nil || expiration.Equal(zeroTime) {
		return -1
	}
	return time.Until(*expiration)
}

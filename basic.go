package cooldown

import (
	"sync/atomic"
	"time"
)

// Basic represents very basic cooldown. It is 'lazy'.
// You can create it simpy with: new(cooldown.Basic)
type Basic struct {
	// expiration is pointer to expiration time.
	// It is nil either if cooldown was just created or if it was reset.
	expiration atomic.Pointer[time.Time]
}

// Set updates state of the cooldown.
func (cooldown *Basic) Set(dur time.Duration) {
	now := time.Now()
	expiration := now.Add(dur)
	// Storing expiration date as pointer
	cooldown.expiration.Store(&expiration)
}

var zeroTime = time.Time{}

// Reset resets the cooldown expiration.
func (cooldown *Basic) Reset() {
	// Clear the expiration date
	// Our structure can handle nil values, so there are no problems
	cooldown.expiration.Store(nil)
}

// Active returns true if cooldown is currently active.
func (cooldown *Basic) Active() bool {
	_, ok := cooldown.active()
	return ok
}

func (cooldown *Basic) active() (time.Time, bool) {
	expiration := cooldown.expiration.Load()
	if expiration == nil {
		return zeroTime, false
	}
	e := *expiration
	if e.Equal(zeroTime) {
		return e, false
	}
	now := time.Now()
	// If expiration date is before current date, it is expired. If it is not,
	// expiration date haven't been passed
	return e, !e.Before(now)
}

// Remaining returns the duration until cooldown expiration.
func (cooldown *Basic) Remaining() time.Duration {
	e, ok := cooldown.active()
	if !ok {
		return 0
	}
	if e.Equal(zeroTime) {
		return -1
	}
	return time.Until(e)
}

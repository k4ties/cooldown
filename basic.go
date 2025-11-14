package cooldown

import (
	"sync"
	"time"
)

// Basic represents very basic cooldown. It is 'lazy'.
// You can create it simpy with: new(cooldown.Basic)
// You can use your own mutex by calling only 'Unsafe' methods with locking
// by yourself.
type Basic struct {
	L sync.Mutex
	// expiration is pointer to expiration time.
	// It is nil either if cooldown was just created or if it was reset.
	expiration,
	// pausedAt is pointer to time when cooldown was paused.
	pausedAt *time.Time
}

// Set updates state of the cooldown. It resets cooldown each call.
// If provided duration is negative, cooldown will just reset.
func (cooldown *Basic) Set(dur time.Duration) {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	cooldown.SetUnsafe(dur)
}

func (cooldown *Basic) SetUnsafe(dur time.Duration) {
	// Reset the cooldown data before setting new one
	// That is mostly important for previously paused cooldowns
	cooldown.ResetUnsafe()
	if dur <= 0 {
		// Nothing to do, cooldown did reset already
		return
	}
	// Storing expiration date pointer
	expiration := time.Now().Add(dur)
	cooldown.expiration = &expiration
}

// Pause pauses the cooldown, if it is NOT already paused.
// It returns true, if cooldown was successfully paused.
func (cooldown *Basic) Pause() bool {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.PauseUnsafe()
}

func (cooldown *Basic) PauseUnsafe() bool {
	state := cooldown.StateUnsafe()
	if !state.Active || state.Paused {
		return false
	}
	pausedAt := time.Now()
	// Store the paused date pointer
	cooldown.pausedAt = &pausedAt
	return true
}

// Resume tries to resume the cooldown, if it is paused.
// It returns true, if cooldown was successfully resumed.
func (cooldown *Basic) Resume() bool {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.ResumeUnsafe()
}

func (cooldown *Basic) ResumeUnsafe() bool {
	if _, ok := cooldown.pausedDateUnsafe(); !ok {
		// Cooldown is not paused.
		return false
	}
	// Store nil pointer to paused date
	cooldown.pausedAt = nil
	return true
}

// TogglePause toggles the pause state of the cooldown.
// It returns true, if cooldown was paused.
func (cooldown *Basic) TogglePause() (paused bool) {
	if cooldown.PausedUnsafe() {
		// Cooldown is paused, so resuming it
		cooldown.ResumeUnsafe()
		return false
	}
	// Cooldown is not paused, so pausing it
	return cooldown.PauseUnsafe()
}

// Paused returns true, if cooldown is currently paused.
func (cooldown *Basic) Paused() bool {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.PausedUnsafe()
}

func (cooldown *Basic) PausedUnsafe() bool {
	_, ok := cooldown.pausedDateUnsafe()
	return ok
}

// Reset resets the cooldown data.
func (cooldown *Basic) Reset() {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	cooldown.ResetUnsafe()
}

func (cooldown *Basic) ResetUnsafe() {
	// Our structure can handle nil values, so there are no problems with
	// storing time pointers as nil
	cooldown.expiration = nil
	cooldown.pausedAt = nil
}

// Active returns true if cooldown is currently active.
func (cooldown *Basic) Active() bool {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.ActiveUnsafe()
}

func (cooldown *Basic) ActiveUnsafe() bool {
	return cooldown.StateUnsafe().Active
}

// Remaining returns the duration until cooldown expiration.
func (cooldown *Basic) Remaining() time.Duration {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.RemainingUnsafe()
}

func (cooldown *Basic) RemainingUnsafe() time.Duration {
	res := cooldown.StateUnsafe()
	if !res.Active {
		return 0
	}
	if res.Paused {
		// Calculate remaining time from paused date
		return res.Expiration.Sub(res.PausedDate)
	}
	// Calculate remaining from current time
	// Note: Expiration can't be zero here
	return time.Until(res.Expiration)
}

func (cooldown *Basic) pausedDateUnsafe() (_ time.Time, _ bool) {
	tPtr := cooldown.pausedAt
	if tPtr == nil {
		return
	}
	t := *tPtr
	return t, !t.IsZero()
}

// BasicState represents the state of the basic cooldown.
type BasicState struct {
	// Active is true if cooldown is active.
	Active,
	// Paused is true if cooldown is paused.
	Paused bool
	// Expiration is the expiration date of the cooldown.
	// Will be zero, if cooldown is not active.
	Expiration,
	// PausedDate is the date when cooldown was paused.
	// Will be zero, if cooldown is not paused.
	PausedDate time.Time
}

// State returns the current state of the basic cooldown.
func (cooldown *Basic) State() (res BasicState) {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.StateUnsafe()
}

func (cooldown *Basic) StateUnsafe() (res BasicState) {
	pausedDate, ok := cooldown.pausedDateUnsafe()
	if ok {
		// Update the result
		res.Paused = true
		res.PausedDate = pausedDate
	}
	expiration := cooldown.expiration
	if expiration == nil {
		// Not active and not paused
		return
	}
	res.Expiration = *expiration
	if res.Expiration.IsZero() {
		cooldown.expiration = nil
		return
	}
	reference := time.Now()
	if res.Paused {
		// Calculate if active from the paused date
		reference = pausedDate
	}
	// If expiration date is before current date, it is expired. If it is not,
	// expiration date haven't been passed
	res.Active = !res.Expiration.Before(reference)
	return
}

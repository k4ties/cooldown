package cooldown

import (
	"sync"
	"time"
)

// Basic represents very basic 'lazy' cooldown.
// You can create it simpy via new(cooldown.Basic).
type Basic struct {
	L sync.RWMutex
	// expiration is time when cooldown expires.
	expiration,
	// pausedAt is time when cooldown was paused.
	pausedAt time.Time
}

// Set updates state of the cooldown.
// If provided duration is negative, cooldown will just reset.
func (cooldown *Basic) Set(dur time.Duration) {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	cooldown.SetUnsafe(dur)
}

func (cooldown *Basic) SetUnsafe(dur time.Duration) {
	cooldown.ResetUnsafe()
	if dur <= 0 {
		return
	}
	cooldown.expiration = time.Now().Add(dur)
}

// Pause pauses cooldown if it is not already paused.
// Returns true if successfully paused.
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
	cooldown.pausedAt = time.Now()
	return true
}

// Resume resumes the cooldown if it is paused.
// Returns true if successfully resumed.
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
	cooldown.pausedAt = time.Time{}
	return true
}

// TogglePause toggles the pause state of the cooldown.
// Returns true if cooldown was paused.
func (cooldown *Basic) TogglePause() (paused bool) {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.TogglePauseUnsafe()
}

func (cooldown *Basic) TogglePauseUnsafe() (paused bool) {
	if cooldown.PausedUnsafe() {
		return cooldown.ResumeUnsafe()
	}
	return cooldown.PauseUnsafe()
}

// Paused returns true if cooldown is paused.
func (cooldown *Basic) Paused() bool {
	cooldown.L.RLock()
	defer cooldown.L.RUnlock()
	return cooldown.PausedUnsafe()
}

func (cooldown *Basic) PausedUnsafe() bool {
	_, ok := cooldown.pausedDateUnsafe()
	return ok
}

// Reset resets the cooldown state.
func (cooldown *Basic) Reset() {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	cooldown.ResetUnsafe()
}

func (cooldown *Basic) ResetUnsafe() {
	cooldown.expiration, cooldown.pausedAt = time.Time{}, time.Time{}
}

// Active returns true if cooldown is currently active.
func (cooldown *Basic) Active() bool {
	cooldown.L.RLock()
	defer cooldown.L.RUnlock()
	return cooldown.ActiveUnsafe()
}

func (cooldown *Basic) ActiveUnsafe() bool {
	return cooldown.StateUnsafe().Active
}

// Remaining returns duration until cooldown expiration.
func (cooldown *Basic) Remaining() time.Duration {
	cooldown.L.RLock()
	defer cooldown.L.RUnlock()
	return cooldown.RemainingUnsafe()
}

func (cooldown *Basic) RemainingUnsafe() time.Duration {
	res := cooldown.StateUnsafe()
	if !res.Active {
		return 0
	}
	if res.Paused {
		return res.Expiration.Sub(res.PausedDate)
	}
	// Note: Expiration can't be zero here
	return time.Until(res.Expiration)
}

func (cooldown *Basic) pausedDateUnsafe() (_ time.Time, _ bool) {
	return cooldown.pausedAt, !cooldown.pausedAt.IsZero()
}

// BasicState represents the state of Basic cooldown.
type BasicState struct {
	// Active marks if cooldown is active.
	Active,
	// Paused marks if cooldown is paused.
	Paused bool
	// Expiration is the expiration date of the cooldown.
	// If cooldown is inactive, it'll be zero time.Time,
	Expiration,
	// PausedDate is the date when cooldown was paused.
	// If it wasn't, it'll be zero time.Time.
	PausedDate time.Time
}

// State returns the current state of the basic cooldown.
func (cooldown *Basic) State() (res BasicState) {
	cooldown.L.Lock()
	defer cooldown.L.Unlock()
	return cooldown.StateUnsafe()
}

func (cooldown *Basic) StateUnsafe() (state BasicState) {
	if pausedDate, ok := cooldown.pausedDateUnsafe(); ok {
		state.Paused = true
		state.PausedDate = pausedDate
	}
	expiration := cooldown.expiration
	if expiration.IsZero() {
		return
	}
	state.Expiration = expiration
	reference := time.Now()
	if state.Paused {
		reference = state.PausedDate
	}
	state.Active = !state.Expiration.Before(reference)
	return state
}

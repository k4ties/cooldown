package cooldown

import "errors"

// StopCause is used to identify reason of cooldown stop.
type StopCause = error

var (
	// ErrStopCauseExpired used when cooldown is expired.
	ErrStopCauseExpired = errors.New("cooldown expired")
	// ErrStopCauseCancelled used when cooldown is cancelled.
	ErrStopCauseCancelled = errors.New("cooldown cancelled")
)

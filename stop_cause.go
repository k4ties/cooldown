package cooldown

import "errors"

// StopCause is used to identify reason of cooldown stop.
type StopCause error

var (
	// StopCauseExpired used when cooldown is expired.
	StopCauseExpired StopCause = errors.New("cooldown expired")
	// StopCauseCancelled used when cooldown is cancelled.
	StopCauseCancelled StopCause = errors.New("cooldown cancelled")
)

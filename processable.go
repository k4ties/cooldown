package cooldown

import "time"

// Processable ...
type Processable interface {
	// Remaining ...
	Remaining() time.Duration
	// UnsafeStop should stop the processable object unsafely (within mutex lock).
	UnsafeStop(cause StopCause, handle bool)
}

// getExpiration ...
func getExpiration(p Processable) time.Time {
	return time.Now().Add(p.Remaining())
}

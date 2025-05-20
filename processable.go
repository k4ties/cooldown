package cooldown

import "time"

// TODO export this and stop method
// processable ...
type processable interface {
	// Remaining ...
	Remaining() time.Duration
	// stop should stop the processable object.
	stop(cause StopCause, handle bool)
}

// getExpiration ...
func getExpiration(p processable) time.Time {
	return time.Now().Add(p.Remaining())
}

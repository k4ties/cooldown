package cooldown

import "errors"

// StopCause is used to identify reason of cooldown stop.
// It extends default error interface.
type StopCause = error

var (
	// ErrStopCauseExpired used when cooldown is expired.
	ErrStopCauseExpired = errors.New("cooldown expired")
	// ErrStopCauseCancelled used when cooldown is canceled.
	ErrStopCauseCancelled = errors.New("cooldown cancelled")
	// ErrStopCauseReset used when cooldown is reset by user.
	ErrStopCauseReset = errors.New("cooldown did reset") // TODO
)

var (
	// ErrStartTimerNotNil is error, that is used, when trying to start
	// cooldown while underlying time.Timer (time.AfterFunc) is not nil.
	ErrStartTimerNotNil = errors.New("cooldown: trying to start cooldown (that is logically inactive), but timer is not nil (either race or bug)")
	// ErrStopTimerNil is error, that is used, when trying to stop cooldown,
	// while underlying time.Timer (time.AfterFunc) is nil. Since we need to
	// prevent cooldown to expire, we need to stop it until its expiration
	// correctly.
	ErrStopTimerNil = errors.New("cooldown: tried to stop (that is logically active) cooldown while timer is nil")
	// ErrRenewTimerNil is error, that is used, when trying to renew cooldown,
	// while underlying time.Timer (time.AfterFunc) is nil. We need to prevent
	// cooldown to expire earlier, than it should (because we're renewing the
	// cooldown).
	ErrRenewTimerNil = errors.New("cooldown: tried to renew cooldown (that is logically active) while timer is nil")
)

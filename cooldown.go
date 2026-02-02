package cooldown

import "time"

// CoolDown is the same Valued cooldown but without values.
// It uses underlying Valued with empty struct as value.
// The handler type is also different here, but its logic is the same, because
// it's just a wrapper for ValuedHandler.
type CoolDown struct {
	valued *Valued[struct{}]
}

// New creates new CoolDown instance.
func New(opts ...Option) *CoolDown {
	cd := &CoolDown{valued: NewValued[struct{}]()} // Don't type options here, they will be directly applied to CoolDown
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cd)
	}
	return cd
}

/*
   The following methods are wrappers for underlying
   Valued methods, so they aren't documented.
   See Valued for more information.
*/

// Renew ...
func (cooldown *CoolDown) Renew() {
	cooldown.valued.Renew(zeroStruct)
}

func (cooldown *CoolDown) RenewUnsafe() {
	cooldown.valued.RenewUnsafe(zeroStruct)
}

// Start ...
func (cooldown *CoolDown) Start(dur time.Duration) {
	cooldown.valued.Start(dur, zeroStruct)
}

func (cooldown *CoolDown) StartUnsafe(dur time.Duration) {
	cooldown.valued.StartUnsafe(dur, zeroStruct)
}

// Stop ...
func (cooldown *CoolDown) Stop() {
	cooldown.valued.Stop(zeroStruct)
}

func (cooldown *CoolDown) StopUnsafe() {
	cooldown.valued.StopUnsafe(zeroStruct)
}

// Pause ...
func (cooldown *CoolDown) Pause() bool {
	return cooldown.valued.Pause(struct{}{})
}

func (cooldown *CoolDown) PauseUnsafe() bool {
	return cooldown.valued.PauseUnsafe(struct{}{})
}

// Resume ...
func (cooldown *CoolDown) Resume() bool {
	return cooldown.valued.Resume(struct{}{})
}

func (cooldown *CoolDown) ResumeUnsafe() bool {
	return cooldown.valued.ResumeUnsafe(struct{}{})
}

// TogglePause ...
func (cooldown *CoolDown) TogglePause() bool {
	return cooldown.valued.TogglePause(struct{}{})
}

func (cooldown *CoolDown) TogglePauseUnsafe() bool {
	return cooldown.valued.TogglePause(struct{}{})
}

// Handler ...
func (cooldown *CoolDown) Handler() Handler {
	return convertFromValuedHandler(cooldown.valued.Handler(), cooldown.valued)
}

// Handle ...
func (cooldown *CoolDown) Handle(handler Handler) {
	cooldown.valued.Handle(convertToValuedHandler(handler, cooldown))
}

// Active ...
func (cooldown *CoolDown) Active() bool {
	return cooldown.valued.Active()
}

func (cooldown *CoolDown) ActiveUnsafe() bool {
	return cooldown.valued.ActiveUnsafe()
}

// Remaining ...
func (cooldown *CoolDown) Remaining() time.Duration {
	return cooldown.valued.Remaining()
}

func (cooldown *CoolDown) RemainingUnsafe() time.Duration {
	return cooldown.valued.RemainingUnsafe()
}

// Valued ...
func (cooldown *CoolDown) Valued() *Valued[struct{}] {
	return cooldown.valued
}

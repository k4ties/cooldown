package cooldown

import "time"

// CoolDown is the same Valued cooldown but without values.
type CoolDown struct {
	valued *Valued[struct{}]
}

// New creates new CoolDown.
func New(h Handler) *CoolDown {
	cd := new(CoolDown)
	cd.valued = NewValued[struct{}](
		convertToValuedHandler(h, cd),
	)
	return cd
}

// Renew ...
func (cooldown *CoolDown) Renew() {
	cooldown.valued.Renew(struct{}{})
}

// Start ...
func (cooldown *CoolDown) Start(dur time.Duration) {
	cooldown.valued.Start(dur, struct{}{})
}

// Stop ...
func (cooldown *CoolDown) Stop() {
	cooldown.valued.Stop(struct{}{})
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

// Remaining ...
func (cooldown *CoolDown) Remaining() time.Duration {
	return cooldown.valued.Remaining()
}

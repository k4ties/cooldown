package cooldown

import (
	"errors"
	"github.com/df-mc/dragonfly/server/event"
)

// Context ...
type Context = event.Context[*CoolDown]

// Handler is handler of the main CoolDown actions.
type Handler interface {
	// HandleStart handles start with the ability to cancel it.
	HandleStart(ctx *Context)
	// HandleRenew handles renew with the ability to cancel it.
	HandleRenew(ctx *Context)
	// HandleTick handles every tick of the cooldown (variable TicksPerSecond). We can handle every
	// second, two seconds or any amount that we want by [tickCount % tps == 0] logic.
	HandleTick(cooldown *CoolDown, current int64)
	// HandleStop handles stop of the CoolDown, with the specified cause.
	HandleStop(cooldown *CoolDown, cause StopCause)
}

// NopHandler is no-operation handler of Handler.
type NopHandler struct{}

func (NopHandler) HandleStart(*Context)            {}
func (NopHandler) HandleRenew(*Context)            {}
func (NopHandler) HandleTick(*CoolDown, int64)     {}
func (NopHandler) HandleStop(*CoolDown, StopCause) {}

// StopCause ...
type StopCause error

var (
	// StopCauseExpired used when cooldown is expired.
	StopCauseExpired StopCause = errors.New("expired")
	// StopCauseCancelled used when cooldown is cancelled.
	StopCauseCancelled StopCause = errors.New("cancelled")
)

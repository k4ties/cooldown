package cooldown

import (
	"errors"
	"github.com/k4ties/cooldown/internal/event"
)

// ContextWithVal ...
type ContextWithVal[T any] = event.Context[*WithVal[T]]

// HandlerWithVal is interface that implements basic handler of the CoolDown actions.
type HandlerWithVal[T any] interface {
	// HandleStart handles start with the ability to cancel it.
	HandleStart(ctx *ContextWithVal[T], val T)
	// HandleRenew handles renew with the ability to cancel it.
	HandleRenew(ctx *ContextWithVal[T], val T)
	// HandleTick handles every tick of the cooldown (variable TicksPerSecond). We can handle every
	// second, two seconds or any amount that we want by [tickCount % tps == 0] logic.
	HandleTick(cooldown *WithVal[T], current int64, val T)
	// HandleStop handles stop of the CoolDown, with the specified cause.
	HandleStop(cooldown *WithVal[T], cause StopCause, val T)
}

// NopHandler is no-operation implementation of Handler.
type NopHandler[T any] struct{}

func (NopHandler[T]) HandleStart(*ContextWithVal[T], T)    {}
func (NopHandler[T]) HandleRenew(*ContextWithVal[T], T)    {}
func (NopHandler[T]) HandleTick(*WithVal[T], int64, T)     {}
func (NopHandler[T]) HandleStop(*WithVal[T], StopCause, T) {}

// StopCause ...
type StopCause error

var (
	// StopCauseExpired used when cooldown is expired.
	StopCauseExpired StopCause = errors.New("expired")
	// StopCauseCancelled used when cooldown is cancelled.
	StopCauseCancelled StopCause = errors.New("cancelled")
)

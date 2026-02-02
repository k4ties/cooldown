package cooldown

import (
	"time"

	"github.com/k4ties/cooldown/internal/event"
)

type ValuedContext[T any] = event.Context[*Valued[T]]

// ValuedHandler allows to handle actions with Valued, additionally providing
// a ValuedContext allowing to cancel the event.
//
// Note: you're NOT allowed to call locking Valued methods on handler events,
// because it is already in lock. Otherwise, it'll cause deadlock. Note, that
// you're still able to use Unsafe methods.
type ValuedHandler[T any] interface {
	// HandleStart handles start of the cooldown allowing user to cancel it via
	// context.
	HandleStart(ctx *ValuedContext[T], dur time.Duration, val T)
	// HandleRenew handles renew allowing user to cancel it via context.
	HandleRenew(ctx *ValuedContext[T], dur time.Duration, val T)
	// HandleStop handles stop of the cooldown. You can identify stop cause by
	// errors.Is method. Example:
	//
	// switch cause {
	// case cooldown.ErrStopCauseExpired:
	//    // ...
	// case cooldown.ErrStopCauseCancelled:
	//    // ...
	// }
	HandleStop(cooldown *Valued[T], cause StopCause, val T)
	// HandlePause handles user pausing the cooldown allowing to cancel event
	// via context.
	HandlePause(ctx *ValuedContext[T], val T)
	// HandleResume handles user resuming the cooldown allowing to cancel event
	// via context.
	HandleResume(ctx *ValuedContext[T], val T)
}

type Context = event.Context[*CoolDown]

// Handler allows to handle actions with CoolDown, additionally providing a
// Context allowing to cancel the event.
//
// Note: you're NOT allowed to call locking CoolDown methods on handler events,
// because it is already in lock. Otherwise, it'll cause deadlock. Note, that
// you're still able to use Unsafe methods.
type Handler interface {
	// HandleStart handles start of the cooldown allowing user to cancel it via
	// context.
	HandleStart(ctx *Context, dur time.Duration)
	// HandleRenew handles cooldown renew allowing user to cancel it via
	// context.
	HandleRenew(ctx *Context, dur time.Duration)
	// HandleStop handles stop of the cooldown. You can identify stop cause by
	// errors.Is method. Example:
	//
	// switch cause {
	// case cooldown.ErrStopCauseExpired:
	//    // ...
	// case cooldown.ErrStopCauseCancelled:
	//    // ...
	// }
	HandleStop(cooldown *CoolDown, cause StopCause)
	// HandlePause handles user pausing the cooldown allowing to cancel event
	// via context.
	HandlePause(ctx *Context)
	// HandleResume handles user resuming the cooldown allowing to cancel event
	// via context.
	HandleResume(ctx *Context)
}

// NopValuedHandler is no-operation implementation of ValuedHandler.
type NopValuedHandler[T any] struct{}

func (NopValuedHandler[T]) HandleStart(*ValuedContext[T], time.Duration, T) {}
func (NopValuedHandler[T]) HandleRenew(*ValuedContext[T], time.Duration, T) {}
func (NopValuedHandler[T]) HandleStop(*Valued[T], StopCause, T)             {}
func (NopValuedHandler[T]) HandlePause(*ValuedContext[T], T)                {}
func (NopValuedHandler[T]) HandleResume(*ValuedContext[T], T)               {}

// NopHandler is no-operation implementation of Handler.
type NopHandler struct{}

func (NopHandler) HandleStart(*Context, time.Duration) {}
func (NopHandler) HandleRenew(*Context, time.Duration) {}
func (NopHandler) HandleStop(*CoolDown, StopCause)     {}
func (NopHandler) HandlePause(*Context)                {}
func (NopHandler) HandleResume(*Context)               {}

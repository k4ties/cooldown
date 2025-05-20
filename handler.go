package cooldown

import "github.com/k4ties/cooldown/internal/event"

type (
	// ValuedContext ...
	ValuedContext[T any] = event.Context[*Valued[T]]
	// Context ...
	Context = event.Context[*CoolDown]
)

// ValuedHandler is interface that implements handler of the all Valued cooldown actions.
type ValuedHandler[T any] interface {
	// HandleStart handles start with the ability to cancel it.
	HandleStart(ctx *ValuedContext[T], val T)
	// HandleRenew handles renew with the ability to cancel it.
	HandleRenew(ctx *ValuedContext[T], val T)
	// HandleStop handles stop of the CoolDown, with the specified cause.
	HandleStop(cooldown *Valued[T], cause StopCause, val T)
}

// Handler is interface that implements handler of the all CoolDown actions.
type Handler interface {
	// HandleStart handles start with the ability to cancel it.
	HandleStart(ctx *Context)
	// HandleRenew handles renew with the ability to cancel it.
	HandleRenew(ctx *Context)
	// HandleStop handles stop of the CoolDown, with the specified cause.
	HandleStop(cooldown *CoolDown, cause StopCause)
}

// NopValuedHandler is no-operation implementation of ValuedHandler.
type NopValuedHandler[T any] struct{}

// NopHandler is no-operation implementation of Handler.
type NopHandler struct{}

func (NopValuedHandler[T]) HandleStart(*ValuedContext[T], T)    {}
func (NopValuedHandler[T]) HandleRenew(*ValuedContext[T], T)    {}
func (NopValuedHandler[T]) HandleStop(*Valued[T], StopCause, T) {}

func (NopHandler) HandleStart(*Context)            {}
func (NopHandler) HandleRenew(*Context)            {}
func (NopHandler) HandleStop(*CoolDown, StopCause) {}

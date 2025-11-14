package cooldown

// ValuedHandler is interface that implements handler of the all Valued
// cooldown actions.
type ValuedHandler[T any] interface {
	// HandleStart handles start with the ability to cancel it via provided
	// context.
	HandleStart(ctx *ValuedContext[T], val T)
	// HandleRenew handles renew with the ability to cancel it via provided
	// context.
	HandleRenew(ctx *ValuedContext[T], val T)
	// HandleStop handles stop of the cooldown. You can check stop cause by
	// errors.Is method. Example:
	// switch cause {
	// case cooldown.ErrStopCauseExpired:
	//    // ...
	// case cooldown.ErrStopCauseCancelled:
	//    // ...
	// }
	HandleStop(cooldown *Valued[T], cause StopCause, val T)
}

// Handler is interface that implements handler of the all CoolDown actions.
type Handler interface {
	// HandleStart handles start with the ability to cancel it via provided
	// context.
	HandleStart(ctx *Context)
	// HandleRenew handles renew with the ability to cancel it via provided
	// context.
	HandleRenew(ctx *Context)
	// HandleStop handles stop of the cooldown. You can check stop cause by
	// errors.Is method. Example:
	// switch cause {
	// case cooldown.ErrStopCauseExpired:
	//    // ...
	// case cooldown.ErrStopCauseCancelled:
	//    // ...
	// }
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

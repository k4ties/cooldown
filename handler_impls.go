package cooldown

import "github.com/k4ties/cooldown/internal/event"

// convertToValuedHandler converts Handler to ValuedHandler[struct{}].
func convertToValuedHandler(from Handler, cd *CoolDown) ValuedHandler[struct{}] {
	return valuedHandler[struct{}]{
		parent:   from,
		cooldown: cd,
	}
}

// convertFromValuedHandler converts valuedHandler[struct{}] to Handler.
func convertFromValuedHandler(from ValuedHandler[struct{}], cd *Valued[struct{}]) Handler {
	return handler{
		parent:   from,
		cooldown: cd,
	}
}

// handler is Handler implementation that redirects actions to parent ValuedHandler.
type handler struct {
	parent   ValuedHandler[struct{}]
	cooldown *Valued[struct{}]
}

// HandleStart ...
func (h handler) HandleStart(parent *Context) {
	ctx := convertContext(event.C(h.cooldown))
	if h.parent.HandleStart(ctx, struct{}{}); ctx.Cancelled() {
		parent.Cancel()
	}
}

// HandleRenew ...
func (h handler) HandleRenew(parent *Context) {
	ctx := convertContext(event.C(h.cooldown))
	if h.parent.HandleRenew(ctx, struct{}{}); ctx.Cancelled() {
		parent.Cancel()
	}
}

// HandleStop ...
func (h handler) HandleStop(_ *CoolDown, cause StopCause) {
	h.parent.HandleStop(h.cooldown, cause, struct{}{})
}

// valuedHandler is ValuedHandler implementation that redirects actions to parent Handler.
type valuedHandler[T any] struct {
	parent   Handler
	cooldown *CoolDown
}

// HandleStart ...
func (handler valuedHandler[T]) HandleStart(parent *ValuedContext[T], _ T) {
	ctx := event.C(handler.cooldown)
	if handler.parent.HandleStart(ctx); ctx.Cancelled() {
		parent.Cancel()
	}
}

// HandleRenew ...
func (handler valuedHandler[T]) HandleRenew(parent *ValuedContext[T], _ T) {
	ctx := event.C(handler.cooldown)
	if handler.parent.HandleRenew(ctx); ctx.Cancelled() {
		parent.Cancel()
	}
}

// HandleStop ...
func (handler valuedHandler[T]) HandleStop(_ *Valued[T], cause StopCause, _ T) {
	handler.parent.HandleStop(handler.cooldown, cause)
}

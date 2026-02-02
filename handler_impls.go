package cooldown

import (
	"time"

	"github.com/k4ties/cooldown/internal/event"
)

func convertToValuedHandler(from Handler, cd *CoolDown) ValuedHandler[struct{}] {
	if from == nil {
		return NopValuedHandler[struct{}]{}
	}
	return valuedHandler[struct{}]{
		parent:   from,
		cooldown: cd,
	}
}

func convertFromValuedHandler(from ValuedHandler[struct{}], cd *Valued[struct{}]) Handler { // todo: is this even required
	if from == nil {
		return NopHandler{}
	}
	return &handler{
		parent:   from,
		cooldown: cd,
	}
}

type handler struct {
	parent   ValuedHandler[struct{}]
	cooldown *Valued[struct{}]
}

var zeroStruct = struct{}{}

func (h *handler) HandleStart(parent *Context, dur time.Duration) {
	ctx := event.C(h.cooldown)
	if h.parent.HandleStart(ctx, dur, zeroStruct); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (h *handler) HandleRenew(parent *Context, dur time.Duration) {
	ctx := event.C(h.cooldown)
	if h.parent.HandleRenew(ctx, dur, zeroStruct); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (h *handler) HandlePause(parent *Context) {
	ctx := event.C(h.cooldown)
	if h.parent.HandlePause(ctx, struct{}{}); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (h *handler) HandleResume(parent *Context) {
	ctx := event.C(h.cooldown)
	if h.parent.HandleResume(ctx, struct{}{}); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (h *handler) HandleStop(_ *CoolDown, cause StopCause) {
	h.parent.HandleStop(h.cooldown, cause, zeroStruct)
}

type valuedHandler[T any] struct {
	parent   Handler
	cooldown *CoolDown
}

func (handler valuedHandler[T]) HandleStart(parent *ValuedContext[T], dur time.Duration, _ T) {
	ctx := event.C(handler.cooldown)
	if handler.parent.HandleStart(ctx, dur); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (handler valuedHandler[T]) HandleRenew(parent *ValuedContext[T], dur time.Duration, _ T) {
	ctx := event.C(handler.cooldown)
	if handler.parent.HandleRenew(ctx, dur); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (handler valuedHandler[T]) HandlePause(parent *ValuedContext[T], _ T) {
	ctx := event.C(handler.cooldown)
	if handler.parent.HandlePause(ctx); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (handler valuedHandler[T]) HandleResume(parent *ValuedContext[T], _ T) {
	ctx := event.C(handler.cooldown)
	if handler.parent.HandleResume(ctx); ctx.Cancelled() {
		parent.Cancel()
	}
}
func (handler valuedHandler[T]) HandleStop(_ *Valued[T], cause StopCause, _ T) {
	handler.parent.HandleStop(handler.cooldown, cause)
}

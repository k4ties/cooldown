package event

func Ctx[T any](v T) *Context[T] {
	return &Context[T]{val: v}
}

type Context[T any] struct {
	cancel bool
	val    T
}

func (ctx *Context[T]) Val() T          { return ctx.val }
func (ctx *Context[T]) Cancelled() bool { return ctx.cancel }
func (ctx *Context[T]) Cancel()         { ctx.cancel = true }

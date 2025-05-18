package cooldown

type Option[T any] = func(*WithVal[T])

func StartFunc[T any](f StartTaskFunc) Option[T] {
	return func(c *WithVal[T]) {
		c.taskFunc = f
	}
}

func WithHandler[T any](h HandlerWithVal[T]) Option[T] {
	return func(w *WithVal[T]) {
		w.Handle(h)
	}
}

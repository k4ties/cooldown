package atomic

import "sync/atomic"

// Value is used to use atomic.Pointer without pointers. (lol)
type Value[T any] struct {
	v *atomic.Pointer[T]
}

// NewValue creates new Value impl.
func NewValue[T any](val ...T) Value[T] {
	v := Value[T]{
		v: &atomic.Pointer[T]{},
	}
	if len(val) > 0 {
		v.Store(val[0])
	}
	return v
}

// Load tries to load a value from parent atomic.Value. If value is nil, zero T and false is
// returned. Then, it tries to type assert the value to T and returns the end result.
func (value Value[T]) Load() (T, bool) {
	v := value.v.Load()
	if v == nil {
		var zero T
		return zero, false
	}

	return *v, true
}

// MustLoad tries to load value. If error was occurred, it'll panic.
func (value Value[T]) MustLoad() T {
	v, ok := value.Load()
	if !ok {
		panic("cannot load value")
	}
	return v
}

// Swap swaps current Value val with the specified one. If currently Value has not val, zero T and
// false is returned.
func (value Value[T]) Swap(val T) (old T, hasOld bool) {
	old, hasOld = value.Load()
	value.Store(val)
	return
}

// CompareAndSwap ...
func (value Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return value.v.CompareAndSwap(&old, &new)
}

// Store ...
func (value Value[T]) Store(val T) {
	value.v.Store(&val)
}

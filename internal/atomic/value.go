package atomic

import "sync/atomic"

// Value is used to use atomic.Value without type assertion and panic risks.
type Value[T any] struct {
	v *atomic.Value
}

// NewValue creates new Value impl.
func NewValue[T any](val ...T) Value[T] {
	v := Value[T]{
		v: &atomic.Value{},
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

	val, ok := v.(T)
	return val, ok
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
	return value.v.CompareAndSwap(old, new)
}

// Store ...
func (value Value[T]) Store(val T) {
	value.v.Store(val)
}

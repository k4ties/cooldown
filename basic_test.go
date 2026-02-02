package cooldown_test

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/k4ties/cooldown"
)

func TestBasic(t *testing.T) {
	b := new(cooldown.Basic)
	testCoolDown(b, cooldownProperties[*cooldown.Basic]{
		reset: func(basic *cooldown.Basic) {
			basic.Reset()
		},
		start: func(basic *cooldown.Basic, duration time.Duration) {
			basic.Set(duration)
		},
		pause: func(basic *cooldown.Basic) bool {
			return basic.Pause()
		},
		resume: func(basic *cooldown.Basic) bool {
			return basic.Resume()
		},
		togglePause: func(basic *cooldown.Basic) bool {
			return basic.TogglePause()
		},
	})(t)
}

type anyCoolDown interface {
	Paused() bool
	Remaining() time.Duration
	Active() bool
}

type cooldownProperties[T anyCoolDown] struct {
	reset         func(T)
	start         func(T, time.Duration)
	pause, resume func(T) bool
	togglePause   func(T) bool
}

func testCoolDown[T anyCoolDown](cd T, props cooldownProperties[T]) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("cooldown functionality test", testCoolDownBasically(cd, props))
		t.Run("pause test", testCoolDownPause(cd, props))
	}
}

func testCoolDownPause[T anyCoolDown](b T, props cooldownProperties[T]) func(*testing.T) {
	return func(t *testing.T) {
		for _, fn := range []func(*testing.T){
			testCoolDownPauseToggle(b, props),
			testCoolDownPauseState(b, props),
		} {
			props.reset(b)
			fn(t)
		}
	}
}

func testCoolDownPauseToggle[T anyCoolDown](cd T, props cooldownProperties[T]) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("TogglePause calls", func(t *testing.T) {
			props.reset(cd)
			props.start(cd, time.Second)

			assert.Equal(t, cd.Paused(), false)
			assert.Equal(t, props.togglePause(cd), true)
			assert.Equal(t, cd.Paused(), true)
			assert.Equal(t, props.togglePause(cd), true)
			assert.Equal(t, cd.Paused(), false)
		})
		t.Run("Pause and Resume methods", func(t *testing.T) {
			props.reset(cd)
			if props.pause(cd) {
				t.Fatal("must not be able to pause cooldown in inactive state")
			}
			props.start(cd, time.Second)

			assert.Equal(t, props.pause(cd), true)
			assert.Equal(t, cd.Paused(), true)
			assert.Equal(t, props.resume(cd), true)
			assert.Equal(t, cd.Paused(), false)

			// Check for calls in unexpected state
			assert.Equal(t, props.resume(cd), false)
			assert.Equal(t, props.pause(cd), true)
			assert.Equal(t, props.pause(cd), false)
			assert.Equal(t, props.resume(cd), true)

			props.reset(cd)
			// Cannot pause while cooldown is not active
			assert.Equal(t, props.pause(cd), false)
		})
		t.Run("Pause() and Resume() methods with delay", func(t *testing.T) {
			props.start(cd, time.Second)
			<-time.After(time.Millisecond)
			assert.Equal(t, props.pause(cd), true)
			<-time.After(time.Millisecond * 10)
			assert.Equal(t, props.resume(cd), true)
		})
	}
}

func testCoolDownPauseState[T anyCoolDown](cd T, props cooldownProperties[T]) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("basic state test", func(t *testing.T) {
			props.start(cd, time.Second)
			assert.Equal(t, cd.Paused(), false)
		})
		t.Run("Set() call while paused", func(t *testing.T) {
			props.start(cd, time.Second)
			assert.Equal(t, props.pause(cd), true)
			assert.Equal(t, cd.Paused(), true)

			props.start(cd, time.Second)
			// The cooldown must NOT be paused now
			assert.Equal(t, cd.Paused(), false)
		})
		t.Run("Active() call while paused", func(t *testing.T) {
			props.start(cd, time.Second)
			assert.Equal(t, props.pause(cd), true)
			assert.Equal(t, cd.Paused(), true)
			// The paused cooldown must be still marked as active
			assert.Equal(t, cd.Active(), true)
		})
	}
}

func testCoolDownBasically[T anyCoolDown](b T, props cooldownProperties[T]) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("active and remaining test", func(t *testing.T) {
			props.reset(b)

			assert.Equal(t, b.Active(), false)
			assert.Equal(t, b.Remaining() <= 0, true)

			props.start(b, time.Second)
			assert.Equal(t, b.Active(), true)

			props.reset(b)
			assert.Equal(t, b.Active(), false)
			assert.Equal(t, b.Remaining() <= 0, true)
		})
		t.Run("negative duration set", func(t *testing.T) {
			props.reset(b)

			props.start(b, 0)
			assert.Equal(t, b.Active(), false)
			props.start(b, -1)
			assert.Equal(t, b.Active(), false)
		})
		t.Run("Remaining() during/after expiration", func(t *testing.T) {
			props.start(b, time.Millisecond)
			<-time.After(5 * time.Millisecond)
			assert.Equal(t, b.Active(), false)
			assert.Equal(t, b.Remaining() <= 0, true)
		})
	}
}

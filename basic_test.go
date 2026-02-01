package cooldown_test

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/k4ties/cooldown"
)

func TestBasic(t *testing.T) {
	b := new(cooldown.Basic)
	testCoolDown(b)(t)
}

type anyCoolDown interface {
	Reset()
	Paused() bool
	Pause() bool
	Resume() bool
	TogglePause() bool
	Remaining() time.Duration
	Active() bool
	Set(time.Duration)
}

func testCoolDown(cd anyCoolDown) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("cooldown functionality test", testCoolDownBasically(cd))
		t.Run("pause test", testCoolDownPause(cd))
	}
}

func testCoolDownPause(b anyCoolDown) func(*testing.T) {
	return func(t *testing.T) {
		for _, fn := range []func(anyCoolDown) func(*testing.T){
			testCoolDownPauseToggle,
			testCoolDownPauseState,
		} {
			b.Reset()
			fn(b)(t)
		}
	}
}

func testCoolDownPauseToggle(cooldown anyCoolDown) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("TogglePause calls", func(t *testing.T) {
			cooldown.Set(time.Second)

			assert.Equal(t, cooldown.Paused(), false)
			assert.Equal(t, cooldown.TogglePause(), true)
			assert.Equal(t, cooldown.Paused(), true)
			assert.Equal(t, cooldown.TogglePause(), false)
			assert.Equal(t, cooldown.Paused(), false)
		})
		t.Run("Pause and Resume methods", func(t *testing.T) {
			cooldown.Set(time.Second)

			assert.Equal(t, cooldown.Pause(), true)
			assert.Equal(t, cooldown.Paused(), true)
			assert.Equal(t, cooldown.Resume(), true)
			assert.Equal(t, cooldown.Paused(), false)

			// Check for calls in unexpected state
			assert.Equal(t, cooldown.Resume(), false)
			assert.Equal(t, cooldown.Pause(), true)
			assert.Equal(t, cooldown.Pause(), false)
			assert.Equal(t, cooldown.Resume(), true)

			cooldown.Reset()
			// Cannot pause while cooldown is not active
			assert.Equal(t, cooldown.Pause(), false)
		})
		t.Run("Pause() and Resume() methods with delay", func(t *testing.T) {
			cooldown.Set(time.Second)
			<-time.After(time.Millisecond)
			assert.Equal(t, cooldown.Pause(), true)
			<-time.After(time.Millisecond * 10)
			assert.Equal(t, cooldown.Resume(), true)
		})
	}
}

func testCoolDownPauseState(cooldown anyCoolDown) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("basic state test", func(t *testing.T) {
			cooldown.Set(time.Second)
			assert.Equal(t, cooldown.Paused(), false)
		})
		t.Run("Set() call while paused", func(t *testing.T) {
			cooldown.Set(time.Second)
			cooldown.Pause()
			assert.Equal(t, cooldown.Paused(), true)

			cooldown.Set(time.Second)
			// The cooldown must NOT be paused now
			assert.Equal(t, cooldown.Paused(), false)
		})
		t.Run("Active() call while paused", func(t *testing.T) {
			cooldown.Set(time.Second)
			cooldown.Pause()
			assert.Equal(t, cooldown.Paused(), true)
			// The paused cooldown must be marked as active
			assert.Equal(t, cooldown.Active(), true)
		})
	}
}

func testCoolDownBasically(b anyCoolDown) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("active and remaining test", func(t *testing.T) {
			b.Reset()

			assert.Equal(t, b.Active(), false)
			assert.Equal(t, b.Remaining() <= 0, true)

			b.Set(time.Second)
			assert.Equal(t, b.Active(), true)

			b.Reset()
			assert.Equal(t, b.Active(), false)
			assert.Equal(t, b.Remaining() <= 0, true)
		})
		t.Run("negative duration set", func(t *testing.T) {
			b.Reset()

			b.Set(0)
			assert.Equal(t, b.Active(), false)
			b.Set(-1)
			assert.Equal(t, b.Active(), false)
		})
		t.Run("Remaining() during/after expiration", func(t *testing.T) {
			b.Set(time.Millisecond)
			<-time.After(5 * time.Millisecond)
			assert.Equal(t, b.Active(), false)
			assert.Equal(t, b.Remaining() <= 0, true)
		})
	}
}

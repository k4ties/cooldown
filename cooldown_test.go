package cooldown_test

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/k4ties/cooldown"
)

func TestCoolDownHandler(t *testing.T) {
	c := cooldown.NewValued[struct{}]()
	// we did not set the handler, so it should return no-op handler instance by default
	assert.Equal(t, c.Handler(), cooldown.NopValuedHandler[struct{}]{})
}

func TestValued(t *testing.T) {
	c := cooldown.NewValued[struct{}]()
	testCoolDown(c, cooldownProperties[*cooldown.Valued[struct{}]]{
		reset: func(c *cooldown.Valued[struct{}]) {
			c.Stop(struct{}{})
		},
		start: func(c *cooldown.Valued[struct{}], duration time.Duration) {
			c.Start(duration, struct{}{})
		},
		pause: func(c *cooldown.Valued[struct{}]) bool {
			return c.Pause(struct{}{})
		},
		resume: func(c *cooldown.Valued[struct{}]) bool {
			return c.Resume(struct{}{})
		},
		togglePause: func(c *cooldown.Valued[struct{}]) bool {
			return c.TogglePause(struct{}{})
		},
	})(t)
}

func TestCoolDownTogglePause(t *testing.T) {
	c := cooldown.NewValued[struct{}]()
	c.Start(time.Second*5, struct{}{})

	for n := range 100 {
		if !c.TogglePause(struct{}{}) {
			t.Fatalf("toggle pause must never return false in active state unless there is a internal error: %d", n)
		}
	}
}

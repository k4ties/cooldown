package cooldown_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/k4ties/cooldown"
)

func TestCoolDownHandler(t *testing.T) {
	c := cooldown.NewValued[struct{}]()
	// we did not set the handler, so it should return no-op handler instance by default
	assert.Equal(t, c.Handler(), cooldown.NopValuedHandler[struct{}]{})
}

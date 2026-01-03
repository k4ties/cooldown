package cooldown_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/k4ties/cooldown"
)

func TestCoolDownHandler(t *testing.T) {
	c := cooldown.NewValued[struct{}]()
	// Returned handler should be no-op
	assert.Equal(t, c.Handler() == cooldown.NopValuedHandler[struct{}]{}, true)
}

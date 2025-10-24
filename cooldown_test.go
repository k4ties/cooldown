package cooldown_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/k4ties/cooldown"
	"testing"
)

func TestCoolDown(t *testing.T) {
	c := cooldown.NewValued[struct{}](nil)
	// Returned handler should be no-op
	assert.Equal(t, c.Handler() == cooldown.NopValuedHandler[struct{}]{}, true)
	// Test basic cooldown
	testBasic(t, c.Basic())
}

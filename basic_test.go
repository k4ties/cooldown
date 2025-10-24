package cooldown_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/k4ties/cooldown"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	testBasic(t, new(cooldown.Basic))
}

func testBasic(t *testing.T, b *cooldown.Basic) {
	assert.Equal(t, b.Active(), false)
	assert.Equal(t, b.Remaining() <= 0, true)

	b.Set(time.Second)
	assert.Equal(t, b.Active(), true)

	b.Reset()
	assert.Equal(t, b.Active(), false)
	assert.Equal(t, b.Remaining() <= 0, true)
}

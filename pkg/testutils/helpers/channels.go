package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// AssertChannelClosed asserts that channel is closed.
func AssertChannelClosed[T any](t *testing.T, ch <-chan T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		assert.Fail(t, "channel is not closed")
	case v, ok := <-ch:
		assert.False(t, ok, "channel is not closed but got value=%v", v)
	}
}

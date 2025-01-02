package retry_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/pkg/retry"
)

func TestOptions_TotalAttempts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  uint
		want uint
	}{
		{
			name: "default",
			got:  0,
			want: 1,
		},
		{
			name: "custom",
			got:  2,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opt := &retry.Options{Retries: tt.got}
			assert.Equal(t, tt.want, opt.TotalAttempts())
		})
	}
}

func TestOptions_GetAttemptDelay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		got     []time.Duration
		attempt uint
		want    time.Duration
	}{
		{
			name:    "first default",
			got:     nil,
			attempt: 0,
			want:    time.Second,
		},
		{
			name:    "second default",
			got:     nil,
			attempt: 1,
			want:    2 * time.Second,
		},
		{
			name:    "first empty",
			got:     []time.Duration{},
			attempt: 0,
			want:    time.Second,
		},
		{
			name:    "second empty",
			got:     []time.Duration{},
			attempt: 0,
			want:    time.Second,
		},
		{
			name:    "first custom",
			got:     []time.Duration{time.Microsecond, time.Millisecond},
			attempt: 0,
			want:    time.Microsecond,
		},
		{
			name:    "second custom",
			got:     []time.Duration{time.Microsecond, time.Millisecond},
			attempt: 1,
			want:    time.Millisecond,
		},
		{
			name:    "more than attempts length",
			got:     []time.Duration{time.Microsecond, time.Millisecond},
			attempt: 2,
			want:    time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opt := &retry.Options{
				Retries:  0,
				Attempts: tt.got,
			}
			assert.Equal(t, tt.want, opt.GetAttemptDelay(tt.attempt))
		})
	}
}

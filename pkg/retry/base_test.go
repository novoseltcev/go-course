package retry_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/retry"
)

var errTest = errors.New("test error")

const testValue = 12

func ExampleDo() {
	attempts := []time.Duration{time.Microsecond, 3 * time.Microsecond, 5 * time.Microsecond}
	retries := uint(len(attempts))

	err := retry.Do(context.Background(), func() error {
		return nil
	}, &retry.Options{
		Attempts: attempts,
		Retries:  retries,
	})

	fmt.Printf("%d retries were made and return error=%T", retries, err)
	// Output:
	// 3 retries were made and return error=<nil>
}

func ExampleDo_error() {
	retries := 0

	err := retry.Do(context.Background(), func() error {
		retries++

		return errTest
	}, &retry.Options{
		Attempts: []time.Duration{time.Microsecond},
		Retries:  3,
	})

	fmt.Printf("%d retries were made and return error=%T", retries, err)

	// Output:
	// 3 retries were made and return error=retry.Error
}

func TestDo_withoutOptions(t *testing.T) {
	t.Parallel()

	err := retry.Do(context.Background(), func() error {
		return nil
	}, nil)

	assert.NoError(t, err)
}

func ExampleDoWithData() {
	retries := 0

	val, err := retry.DoWithData(context.Background(), func() (int, error) {
		retries++

		return testValue, nil
	}, &retry.Options{
		Attempts: []time.Duration{time.Microsecond},
		Retries:  3,
	})

	fmt.Printf("%d retries were made and return value=%d and error=%T", retries, val, err)

	// Output:
	// 1 retries were made and return value=12 and error=<nil>
}

func ExampleDoWithData_error() {
	retries := 0

	val, err := retry.DoWithData(context.Background(), func() (int, error) {
		retries++

		return 0, errTest
	}, &retry.Options{
		Attempts: []time.Duration{time.Microsecond},
		Retries:  3,
	})

	fmt.Printf("%d retries were made and return value=%d and error=%T", retries, val, err)

	// Output:
	// 3 retries were made and return value=0 and error=retry.Error
}

func TestDoWithData_withoutOptions(t *testing.T) {
	t.Parallel()

	data, err := retry.DoWithData(context.Background(), func() (int, error) {
		return testValue, nil
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, testValue, data)
}

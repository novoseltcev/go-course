package retry_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/retry"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

func ExamplePgExec() {
	attempts := []time.Duration{time.Microsecond, 3 * time.Microsecond, 5 * time.Microsecond}
	retries := 0

	err := retry.PgExec(context.Background(), func() error {
		retries++

		return nil
	}, &retry.Options{
		Attempts: attempts,
		Retries:  uint(len(attempts)),
	})

	fmt.Printf("%d retries were made and return error=%T", retries, err)
	// Output:
	// 1 retries were made and return error=<nil>
}

func ExamplePgExec_error() {
	retries := 0

	err := retry.PgExec(context.Background(), func() error {
		retries++

		return testutils.Err
	}, &retry.Options{
		Attempts: []time.Duration{time.Microsecond},
		Retries:  3,
	})

	fmt.Printf("%d retries were made and return error=%T", retries, err)

	// Output:
	// 1 retries were made and return error=retry.Error
}

func TestPgExec_withoutOptions(t *testing.T) {
	t.Parallel()

	err := retry.PgExec(context.Background(), func() error {
		return nil
	}, nil)

	assert.NoError(t, err)
}

func TestPgExec_pgConnError(t *testing.T) {
	t.Parallel()

	retries := 0
	testErr := pgconn.PgError{Code: pgerrcode.ConnectionFailure}

	err := retry.PgExec(context.Background(), func() error {
		retries++

		return &testErr
	}, &retry.Options{Attempts: []time.Duration{time.Microsecond}})

	require.Equal(t, 3, retries)
	require.Error(t, err)
	assert.ErrorIs(t, err, &testErr)
}

func ExamplePgSelect() {
	retries := 0

	val, err := retry.PgSelect(context.Background(), func() (int, error) {
		retries++

		return testutils.INT, nil
	}, &retry.Options{
		Attempts: []time.Duration{time.Microsecond},
		Retries:  3,
	})

	fmt.Printf("%d retries were made and return value=%d and error=%T", retries, val, err)

	// Output:
	// 1 retries were made and return value=10 and error=<nil>
}

func ExamplePgSelect_error() {
	retries := 0

	val, err := retry.PgSelect(context.Background(), func() (int, error) {
		retries++

		return 0, testutils.Err
	}, &retry.Options{
		Attempts: []time.Duration{time.Microsecond},
		Retries:  3,
	})

	fmt.Printf("%d retries were made and return value=%d and error=%T", retries, val, err)

	// Output:
	// 1 retries were made and return value=0 and error=retry.Error
}

func TestPgSelect_withoutOptions(t *testing.T) {
	t.Parallel()

	data, err := retry.PgSelect(context.Background(), func() (int, error) {
		return testutils.INT, nil
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, testutils.INT, data)
}

func TestPgSelect_pgConnError(t *testing.T) {
	t.Parallel()

	retries := 0
	testErr := pgconn.PgError{Code: pgerrcode.ConnectionFailure}

	_, err := retry.PgSelect(context.Background(), func() (int, error) {
		retries++

		return 0, &testErr
	}, &retry.Options{Attempts: []time.Duration{time.Microsecond}})

	require.Equal(t, 3, retries)
	require.Error(t, err)
	assert.ErrorIs(t, err, &testErr)
}

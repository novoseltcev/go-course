package reporters_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/internal/schemas"
)

type clientMock struct {
	res *http.Response
	err error
}

func (c clientMock) Do(_ *http.Request) (*http.Response, error) {
	return c.res, c.err
}

const (
	baseURL   = "http://0.0.0.0:8080"
	secretKey = "secret-key"
)

func TestReportsNil(t *testing.T) {
	t.Parallel()

	reporter := reporters.NewAPIReporter(clientMock{res: &http.Response{}, err: nil}, baseURL, secretKey)

	assert.NoError(t, reporter.Report(context.TODO(), nil))
}

func TestReportsWithoutCheckSum(t *testing.T) {
	t.Parallel()

	reporter := reporters.NewAPIReporter(clientMock{res: &http.Response{}, err: nil}, baseURL, "")

	assert.NoError(t, reporter.Report(context.TODO(), []schemas.Metric{}))
}

var errSome = errors.New("some error")

func TestReportsFailedRequest(t *testing.T) {
	t.Parallel()

	reporter := reporters.NewAPIReporter(clientMock{res: nil, err: errSome}, baseURL, secretKey)

	assert.ErrorIs(t, reporter.Report(context.TODO(), []schemas.Metric{}), errSome)
}

func TestReportsSuccess(t *testing.T) {
	t.Parallel()

	reporter := reporters.NewAPIReporter(clientMock{res: &http.Response{
		Body: io.NopCloser(bytes.NewBufferString("")),
	}, err: nil}, baseURL, secretKey)

	assert.NoError(t, reporter.Report(context.TODO(), []schemas.Metric{}))
}

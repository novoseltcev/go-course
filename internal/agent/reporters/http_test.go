//nolint: paralleltest
package reporters_test

import (
	"context"
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/pkg/retry"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

func TestReportSuccessBase(t *testing.T) {
	value := 10.123
	delta := int64(10)

	tests := []struct {
		name string
		got  []schemas.Metric
		want string
	}{
		{
			name: "empty",
			got:  []schemas.Metric{},
			want: `[]`,
		},
		{
			name: "gauge",
			got: []schemas.Metric{
				{
					ID:    testutils.STRING,
					MType: schemas.Gauge,
					Value: &value,
				},
			},
			want: `[{"id":"string","type":"gauge","value":10.123}]`,
		},
		{
			name: "counter",
			got: []schemas.Metric{
				{
					ID:    testutils.STRING,
					MType: schemas.Counter,
					Delta: &delta,
				},
			},
			want: `[{"id":"string","type":"counter","delta":10}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			gock.New(testutils.URL).Post("/").JSON(tt.want).Reply(http.StatusOK)

			reporter := reporters.NewHTTPClient(http.DefaultClient, testutils.URL)

			assert.NoError(t, reporter.Report(context.TODO(), tt.got))
		})
	}
}

func TestReportSuccessWithCheckSum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer gock.Off()

	hasher := NewMockhasher(ctrl)
	reporter := reporters.NewHTTPClient(http.DefaultClient, testutils.URL, reporters.WithCheckSum(hasher))

	hasher.EXPECT().GetHash([]byte(`[]`)).Return([]byte{1, 2, 3}, nil)
	gock.New(testutils.URL).Post("/").MatchHeader("Hashsha256", hex.Dump([]byte{1, 2, 3}))

	assert.NoError(t, reporter.Report(context.TODO(), []schemas.Metric{}))
}

func TestReportSuccessWithOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer gock.Off()

	hasher := NewMockhasher(ctrl)
	compressor := NewMockcompressor(ctrl)
	reporter := reporters.NewHTTPClient(
		http.DefaultClient,
		testutils.URL,
		reporters.WithCheckSum(hasher),
		reporters.WithCompression(compressor),
	)

	compressor.EXPECT().Compress([]byte{1, 2, 3}).Return([]byte{4, 5}, nil)
	hasher.EXPECT().GetHash([]byte{4, 5}).Return([]byte{6}, nil)
	gock.New(testutils.URL).
		Post("/").
		MatchHeader("Content-Encoding", "gzip").
		MatchHeader("Hashsha256", hex.Dump([]byte{6})).
		JSON([]byte{4, 5})

	assert.NoError(t, reporter.Report(context.TODO(), []schemas.Metric{}))
}

func TestReportSuccessWithRetry(t *testing.T) {
	defer gock.Off()

	gock.New(testutils.URL).Post("/").JSON([]byte(`[]`)).Times(2)

	reporter := reporters.NewHTTPClient(http.DefaultClient, testutils.URL, reporters.WithRetry(retry.Options{
		Attempts: []time.Duration{time.Millisecond},
		Retries:  2,
	}))

	assert.NoError(t, reporter.Report(context.TODO(), []schemas.Metric{}))
}

func TestReportFailedRetries(t *testing.T) {
	defer gock.Off()

	reporter := reporters.NewHTTPClient(http.DefaultClient, testutils.URL, reporters.WithRetry(retry.Options{
		Attempts: []time.Duration{time.Millisecond},
		Retries:  2,
	}))

	gock.New(testutils.URL).Post("/").Times(2).ReplyError(testutils.Err)

	assert.ErrorContains(t,
		reporter.Report(context.TODO(), []schemas.Metric{}),
		"error during send request: All attempts fail:",
	)
}

func TestReportFailedCompression(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	compressor := NewMockcompressor(ctrl)
	reporter := reporters.NewHTTPClient(http.DefaultClient, testutils.URL, reporters.WithCompression(compressor))

	compressor.EXPECT().Compress([]byte(`[]`)).Return(nil, testutils.Err)

	assert.EqualError(t,
		reporter.Report(context.TODO(), []schemas.Metric{}),
		"cannot compress metrics: "+testutils.Err.Error(),
	)
}

func TestReportFailedCheckSum(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	hasher := NewMockhasher(ctrl)
	reporter := reporters.NewHTTPClient(http.DefaultClient, testutils.STRING, reporters.WithCheckSum(hasher))

	hasher.EXPECT().GetHash([]byte(`[]`)).Return(nil, testutils.Err)

	assert.EqualError(t, reporter.Report(context.TODO(), []schemas.Metric{}), "cannot get hash: "+testutils.Err.Error())
}

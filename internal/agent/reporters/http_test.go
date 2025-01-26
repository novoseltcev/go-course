// nolint: paralleltest
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

var (
	value       = float64(testutils.FLAOT)
	delta       = int64(testutils.INT)
	testMetrics = []schemas.Metric{{ID: testutils.STRING, MType: schemas.Counter}}
)

const testMetricsJSON = `[{"id":"string","type":"counter"}]`

func TestHTTPReporter_ReportSuccessBase(t *testing.T) {
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

			reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.URL)

			assert.NoError(t, reporter.Report(context.TODO(), tt.got))
		})
	}
}

func TestHTTPReporter_ReportSuccessWithCheckSum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer gock.Off()

	hasher := NewMockhasher(ctrl)
	reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.URL, reporters.WithCheckSum(hasher))

	hasher.EXPECT().GetHash([]byte(testMetricsJSON)).Return([]byte{1, 2, 3}, nil)
	gock.New(testutils.URL).Post("/").MatchHeader("Hashsha256", hex.Dump([]byte{1, 2, 3}))

	assert.NoError(t, reporter.Report(context.TODO(), testMetrics))
}

func TestHTTPReporter_ReportSuccessWithOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer gock.Off()

	hasher := NewMockhasher(ctrl)
	encryptor := NewMockencryptor(ctrl)
	compressor := NewMockcompressor(ctrl)
	reporter := reporters.NewHTTPReporter(
		http.DefaultClient,
		testutils.URL,
		reporters.WithCheckSum(hasher),
		reporters.WithEncryption(encryptor),
		reporters.WithCompression(compressor),
	)

	encryptor.EXPECT().Encrypt([]byte(testMetricsJSON)).Return([]byte{1, 2, 3}, nil)
	compressor.EXPECT().Compress([]byte{1, 2, 3}).Return([]byte{4, 5}, nil)
	hasher.EXPECT().GetHash([]byte{4, 5}).Return([]byte{6}, nil)
	gock.New(testutils.URL).
		Post("/").
		MatchHeader("Content-Encoding", "gzip").
		MatchHeader("Hashsha256", hex.Dump([]byte{6})).
		JSON([]byte{4, 5})

	assert.NoError(t, reporter.Report(context.TODO(), testMetrics))
}

func TestHTTPReporter_ReportSuccessWithRetry(t *testing.T) {
	defer gock.Off()

	gock.New(testutils.URL).Post("/").JSON([]byte(testMetricsJSON)).Times(2)

	reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.URL, reporters.WithRetry(retry.Options{
		Attempts: []time.Duration{time.Millisecond},
		Retries:  2,
	}))

	assert.NoError(t, reporter.Report(context.TODO(), testMetrics))
}

func TestHTTPReporter_ReportFailedRetries(t *testing.T) {
	defer gock.Off()

	reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.URL, reporters.WithRetry(retry.Options{
		Attempts: []time.Duration{time.Millisecond},
		Retries:  2,
	}))

	gock.New(testutils.URL).Post("/").Times(2).ReplyError(testutils.Err)

	assert.ErrorContains(t,
		reporter.Report(context.TODO(), testMetrics),
		"error during send request: All attempts fail:",
	)
}

func TestHTTPReporter_ReportFailedEncription(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	enc := NewMockencryptor(ctrl)
	reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.URL, reporters.WithEncryption(enc))

	enc.EXPECT().Encrypt([]byte(testMetricsJSON)).Return(nil, testutils.Err)

	assert.EqualError(t,
		reporter.Report(context.TODO(), testMetrics),
		"cannot encrypt metrics: "+testutils.Err.Error(),
	)
}

func TestHTTPReporter_ReportFailedCompression(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	compressor := NewMockcompressor(ctrl)
	reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.URL, reporters.WithCompression(compressor))

	compressor.EXPECT().Compress([]byte(testMetricsJSON)).Return(nil, testutils.Err)

	assert.EqualError(t,
		reporter.Report(context.TODO(), testMetrics),
		"cannot compress metrics: "+testutils.Err.Error(),
	)
}

func TestHTTPReporter_ReportFailedCheckSum(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	hasher := NewMockhasher(ctrl)
	reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.STRING, reporters.WithCheckSum(hasher))

	hasher.EXPECT().GetHash([]byte(testMetricsJSON)).Return(nil, testutils.Err)

	assert.EqualError(t, reporter.Report(context.TODO(), testMetrics), "cannot get hash: "+testutils.Err.Error())
}

func TestHTTPReporter_ReportWithNilRequestBodySkipSend(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	enc := NewMockencryptor(ctrl)
	reporter := reporters.NewHTTPReporter(http.DefaultClient, testutils.URL, reporters.WithEncryption(enc))

	enc.EXPECT().Encrypt(gomock.Any()).Return([]byte{}, nil)
	gock.New(testutils.URL).Post("/").Times(0)

	assert.NoError(t, reporter.Report(context.TODO(), testMetrics))
}

func TestHTTPReporter_ReportFailsCreateRequest(t *testing.T) {
	t.Parallel()

	reporter := reporters.NewHTTPReporter(http.DefaultClient, "\t")

	assert.ErrorContains(t,
		reporter.Report(context.TODO(), testMetrics),
		"cannot create request: ",
	)
}

package workers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/novoseltcev/go-course/internal/agent/workers"
	"github.com/novoseltcev/go-course/internal/schema"
)

type ClientMock struct{}

func (c ClientMock) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{}, nil //nolint:exhaustruct
}

const (
	baseURL   = "http://0.0.0.0:8080"
	secretKey = "secret-key"
)

func TestSendMetrics(t *testing.T) {
	t.Parallel()

	ch := make(chan schema.Metric)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	defer cancel()

	go workers.SendMetrics(ctx, ch, ClientMock{}, 1, baseURL, secretKey)

	value := 123.321
	delta := int64(2)
	ch <- schema.Metric{ID: "Some", MType: schema.Gauge, Value: &value}
	ch <- schema.Metric{ID: "Some", MType: schema.Counter, Delta: &delta}

	close(ch)
}

func BenchmarkSendMetrics(b *testing.B) {
	ch := make(chan schema.Metric)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	defer cancel()

	go workers.SendMetrics(ctx, ch, ClientMock{}, 1, baseURL, secretKey)

	value := 123.321
	delta := int64(2)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ch <- schema.Metric{ID: "Some", MType: schema.Counter, Delta: &delta}
		ch <- schema.Metric{ID: "Some", MType: schema.Gauge, Value: &value}
	}

	<-ctx.Done()
}

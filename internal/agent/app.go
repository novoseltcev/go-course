package agent

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/agent/collectors"
	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/pkg/compress"
	"github.com/novoseltcev/go-course/pkg/cryptoalg"
	"github.com/novoseltcev/go-course/pkg/hash"
	"github.com/novoseltcev/go-course/pkg/retry"
	"github.com/novoseltcev/go-course/pkg/workers"
)

type Agent struct {
	cfg *Config
	r   reporters.Reporter
}

func NewAgent(cfg *Config) *Agent {
	compressor, err := compress.NewGzip(gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}

	opts := []reporters.Option{
		reporters.WithCompression(compressor),
		reporters.WithRetry(retry.Options{
			Retries:  3, //nolint:mnd
			Attempts: []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
		}),
	}

	if cfg.CryptoKey != "" {
		enc, err := cryptoalg.NewPKCS1v15EncryptorFromFile(cfg.CryptoKey)
		if err != nil {
			log.Fatal(err)
		}

		opts = append(opts, reporters.WithEncryption(enc))
	}

	if cfg.SecretKey != "" {
		opts = append(opts, reporters.WithCheckSum(hash.NewHMAC(cfg.SecretKey, sha256.New)))
	}

	return &Agent{
		cfg: cfg,
		r:   reporters.NewHTTPClient(http.DefaultClient, cfg.Address, opts...),
	}
}

func (s *Agent) Start(ctx context.Context) {
	runtimeMetricCh := workers.Producer(ctx, collectors.CollectRuntimeMetrics, s.cfg.PollInterval)
	coreMetricCh := workers.Producer(ctx, collectors.CollectCoreMetrics, s.cfg.PollInterval)

	metricCh := workers.FanIn(ctx, runtimeMetricCh, coreMetricCh)

	go workers.AntiFraudConsumer(ctx, metricCh, s.r.Report, s.cfg.RateLimit)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}

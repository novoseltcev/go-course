package agent

import (
	"context"
	"net/http"
	_ "net/http/pprof" //nolint:gosec

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/agent/workers"
)

type Agent struct {
	config Config
	client http.Client
}

func NewAgent(config Config) *Agent {
	return &Agent{
		config: config,
		client: *http.DefaultClient,
	}
}

func (s *Agent) Start(ctx context.Context) {
	runtimeMetricCh := workers.CollectMetrics(ctx, s.config.PollInterval)
	coreMetricCh := workers.CollectCoreMetrics(ctx, s.config.PollInterval)

	metricCh := workers.FanIn(ctx, runtimeMetricCh, coreMetricCh)

	go workers.SendMetrics(ctx, metricCh, &s.client, s.config.RateLimit, "http://"+s.config.Address, s.config.SecretKey)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}

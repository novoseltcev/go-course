package agent

import (
	"context"
	"net/http"

	"github.com/novoseltcev/go-course/internal/agent/workers"
	_ "net/http/pprof"
)

type Agent struct {
	config Config
	client http.Client
}

func NewAgent(config Config) *Agent {
	return &Agent{
		config: config,
		client: http.Client{},
	}
}

func (s *Agent) Start(ctx context.Context) {
	runtimeMetricCh := workers.CollectMetrics(ctx, s.config.PollInterval)
	coreMetricCh := workers.CollectCoreMetrics(ctx, s.config.PollInterval)

	metricCh := workers.FanIn(ctx, runtimeMetricCh, coreMetricCh)

	go workers.SendMetrics(metricCh, s.config.RateLimit, &s.client, "http://" + s.config.Address, s.config.SecretKey)

	http.ListenAndServe(":9000", nil)
    <- ctx.Done()
}

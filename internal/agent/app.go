package agent

import (
	"context"
	"net/http"

	"github.com/novoseltcev/go-course/internal/agent/workers"
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
	metricCh := workers.CollectMetrics(ctx, s.config.PollInterval)

	go workers.SendMetrics(metricCh, s.config.RateLimit, &s.client, "http://" + s.config.Address, s.config.SecretKey)

    <- ctx.Done()
}

package agent

import (
	"context"
	"net/http"
	_ "net/http/pprof" //nolint:gosec

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/agent/collectors"
	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/pkg/workers"
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
	runtimeMetricCh := workers.Producer(ctx, collectors.CollectRuntimeMetrics, s.config.PollInterval)
	coreMetricCh := workers.Producer(ctx, collectors.CollectCoreMetrics, s.config.PollInterval)

	metricCh := workers.FanIn(ctx, runtimeMetricCh, coreMetricCh)

	reporter := reporters.NewAPIReporter(http.DefaultClient, "http://"+s.config.Address, s.config.SecretKey)
	go workers.AntiFraudConsumer(ctx, metricCh, reporter.Report, s.config.RateLimit)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}

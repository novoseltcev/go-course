package agent

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"

	"github.com/novoseltcev/go-course/internal/agent/workers"
)

type Agent struct {
	config Config
	cron *cron.Cron
	counterStorage map[string]int64
	gaugeStorage map[string]float64
	client http.Client
}

func NewAgent(config Config) *Agent {
	return &Agent{
		config: config,
		cron: cron.New(),
		counterStorage: make(map[string]int64),
		gaugeStorage: make(map[string]float64),
		client: http.Client{},
	}
}

func (s *Agent) Start() {
	s.cron.AddFunc(fmt.Sprintf("@every %ds", s.config.PollInterval), workers.CollectMetrics(&s.counterStorage, &s.gaugeStorage))
	s.cron.AddFunc(fmt.Sprintf("@every %ds", s.config.ReportInterval), workers.SendMetrics(&s.counterStorage, &s.gaugeStorage, &s.client, "http://" + s.config.Address))

	defer s.cron.Stop()
	s.cron.Start()

    quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <-quitChannel
}

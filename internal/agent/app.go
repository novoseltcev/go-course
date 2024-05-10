package agent

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/novoseltcev/go-course/internal/agent/workers"
)

type Agent struct {
	config Config
	counterStorage map[string]int64
	gaugeStorage map[string]float64
	client http.Client
}

func NewAgent(config Config) *Agent {
	return &Agent{
		config: config,
		counterStorage: make(map[string]int64),
		gaugeStorage: make(map[string]float64),
		client: http.Client{},
	}
}

func (s *Agent) Start() {
    quitChannel := make(chan os.Signal, 1)
	defer close(quitChannel)

	go func() {
		fmt.Println("init CollectMetrics worker")
		for {
			workers.CollectMetrics(&s.counterStorage, &s.gaugeStorage)

			select {
			case <-quitChannel:
				fmt.Println("stop CollectMetrics worker")
				return
			default:
				time.Sleep(time.Duration(s.config.PollInterval) * time.Second)
			}
		}
	}()
	
	go func() {
		fmt.Println("init SendMetrics worker")
		for {
			err := workers.SendMetrics(&s.counterStorage, &s.gaugeStorage, &s.client, "http://" + s.config.Address, s.config.SecretKey)
			if err != nil {
				fmt.Println(err.Error())
			}

			select {
			case <-quitChannel:
				fmt.Println("stop SendMetrics worker")
				return
			default:
				time.Sleep(time.Duration(s.config.ReportInterval) * time.Second)
			}
		}
	}()
	
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <- quitChannel
}

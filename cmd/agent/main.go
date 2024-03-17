package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/robfig/cron"
)


func main() {
	client := http.Client{}
	counterStorage := make(agent.Storage[int64])
	gaugeStorage := make(agent.Storage[float64])

	var pollInterval time.Duration = 2
	var reportInterval time.Duration = 10
	var baseURL = "http://0.0.0.0:8080"

	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %ds", pollInterval), agent.CollectMetrics(&counterStorage, &gaugeStorage))
	c.AddFunc(fmt.Sprintf("@every %ds", reportInterval), agent.SendMetrics(&counterStorage, &gaugeStorage, &client, baseURL))

	defer c.Stop()
	c.Start()
	

    quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <-quitChannel
}

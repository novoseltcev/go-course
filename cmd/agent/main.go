package main

import (
	"fmt"
	"net/http"
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
	var baseUrl = "http://localhost:8080"

	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %ds", pollInterval), agent.CollectMetrics(&counterStorage, &gaugeStorage))
	c.AddFunc(fmt.Sprintf("@every %ds", reportInterval), agent.SendMetrics(&counterStorage, &gaugeStorage, &client, baseUrl))
	c.Start()
	defer c.Stop()
	fmt.Scanln()
}

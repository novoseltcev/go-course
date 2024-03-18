package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/novoseltcev/go-course/internal/shared"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"
)

var address *shared.NetAddress
var reportInterval *time.Duration
var pollInterval *time.Duration

func main() {
	parseFlags()

	client := http.Client{}
	counterStorage := make(agent.Storage[int64])
	gaugeStorage := make(agent.Storage[float64])

	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %ds", pollInterval), agent.CollectMetrics(&counterStorage, &gaugeStorage))
	c.AddFunc(fmt.Sprintf("@every %ds", reportInterval), agent.SendMetrics(&counterStorage, &gaugeStorage, &client, "http://" + address.String()))

	defer c.Stop()
	c.Start()
	

    quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <-quitChannel
}

func parseFlags() {
	pflag.Var(address, "a", "Net address host:port")
	reportInterval = pflag.Duration("r", 10, "send metrics to server interval")
	pollInterval = pflag.Duration("p", 2, "poll runtime metrics interval")
	pflag.Parse()
}

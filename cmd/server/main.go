package main

import (
	"net/http"
	"github.com/spf13/pflag"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/shared"
)


var address *shared.NetAddress

func main() {
	parseFlags()

	var counterStorage storage.Storage[storage.Counter] = storage.MemStorage[storage.Counter]{Metrics: make(map[string]storage.Counter)}
	var gaugeStorage storage.Storage[storage.Gauge] = storage.MemStorage[storage.Gauge]{Metrics: make(map[string]storage.Gauge)}

	r := endpoints.MetricRouter(&counterStorage, &gaugeStorage)

	if err := http.ListenAndServe(address.String(), r); err != nil {
		panic(err)
	}
}

func parseFlags() {
	pflag.Var(address, "a", "Net address host:port")
	pflag.Parse()
}

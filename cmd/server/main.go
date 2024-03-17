package main

import (
	"net/http"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func main() {
	var counterStorage storage.Storage[storage.Counter] = storage.MemStorage[storage.Counter]{Metrics: make(map[string]storage.Counter)}
	var gaugeStorage storage.Storage[storage.Gauge] = storage.MemStorage[storage.Gauge]{Metrics: make(map[string]storage.Gauge)}

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, endpoints.UpdateMetric(&counterStorage, &gaugeStorage))

	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}

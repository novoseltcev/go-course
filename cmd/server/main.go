package main

import (
	"net/http"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func main() {
	var counterStorage storage.Storage[int64] = storage.MemStorage[int64]{Metrics: make(map[string]int64)}
	var gaugeStorage storage.Storage[float64] = storage.MemStorage[float64]{Metrics: make(map[string]float64)}

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, endpoints.UpdateMetric(&counterStorage, &gaugeStorage))

	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}

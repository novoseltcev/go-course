package main

import (
	"net/http"
	"strings"
	"strconv"
	"github.com/novoseltcev/go-course/internal/storage"
)


func main() {
	var counterStorage storage.Storage[int64] = storage.MemStorage[int64]{Metrics: make(map[string]int64)}
	var gaugeStorage storage.Storage[float64] = storage.MemStorage[float64]{Metrics: make(map[string]float64)}

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateMetricHandler(&counterStorage, &gaugeStorage))

	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}

func updateMetricHandler(counterStorage *storage.Storage[int64], gaugeStorage *storage.Storage[float64]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var pathParams = strings.Split(
			strings.TrimPrefix(req.URL.Path, `/update/`),
			`/`,
		)

		if len(pathParams) < 3 {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		metricType := pathParams[0]
		metricName := pathParams[1]
		metricValue := pathParams[2]

		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		
			(*gaugeStorage).Update(metricName, value)
		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			(*counterStorage).Update(metricName, value)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

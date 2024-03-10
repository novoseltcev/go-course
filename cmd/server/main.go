package main

import (
	"net/http"
	"strings"
	"strconv"
	"github.com/novoseltcev/go-course/internal/storage"
)


func main() {
	var counterStorage storage.Storage[int64] = storage.MemStorage[int64]{}
	var gaugeStorage storage.Storage[float64] = storage.MemStorage[float64]{}

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateMetricHandler(&counterStorage, &gaugeStorage))

	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}

func updateMetricHandler(counterStorage *storage.Storage[int64], gaugeStorage *storage.Storage[float64]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost || req.Header.Get("content-type") != "text/plain" {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var path_params = strings.Split(
			strings.TrimPrefix(req.URL.Path, `/update/`),
			`/`,
		)

		if len(path_params) < 3 {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		metric_type := path_params[0]
		metric_name := path_params[1]
		metric_value := path_params[2]

		switch metric_type {
		case "gauge":
			value, err := strconv.ParseFloat(metric_value, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		
			(*gaugeStorage).Update(metric_name, value)
		case "counter":
			value, err := strconv.ParseInt(metric_value, 10, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			(*counterStorage).Update(metric_name, value)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

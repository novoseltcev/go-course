package endpoints

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func UpdateMetric(counterStorage *storage.Storage[storage.Counter], gaugeStorage *storage.Storage[storage.Gauge]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("Handle %s\n", req.RequestURI)
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
			
			(*gaugeStorage).Update(metricName, storage.Gauge(value))
		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			(*counterStorage).Update(metricName, storage.Counter(value))
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		
		res.WriteHeader(http.StatusOK)
	}
}

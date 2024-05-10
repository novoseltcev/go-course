package workers

import (
	"net/http"
	"testing"
)

type ClientMock struct {}

func (c ClientMock) Do (*http.Request) (*http.Response, error) {
    return &http.Response{}, nil
}


func TestSendMetrics(t *testing.T) {
	counterStorage :=  map[string]int64{"SomeCounter": 1}
	gaugeStorage :=  map[string]float64{"SomeGauge": 1.0}
	var client Client = ClientMock{}
	baseURL := "http://0.0.0.0:8080"

	SendMetrics(&counterStorage, &gaugeStorage, client, baseURL, "secret-key")
}

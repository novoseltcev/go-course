package workers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/internal/model"
)

type ClientMock struct {}

func (c ClientMock) Do (*http.Request) (*http.Response, error) {
    return &http.Response{}, nil
}


func TestSendMetrics(t *testing.T) {
	counterStorage :=  map[string]model.Counter{"SomeCounter": 1}
	gaugeStorage :=  map[string]model.Gauge{"SomeGauge": 1.0}
	var client Client = ClientMock{}
	baseURL := "http://0.0.0.0:8080"

	SendMetrics(&counterStorage, &gaugeStorage, client, baseURL)()

	assert.Empty(t, counterStorage)
	assert.Empty(t, gaugeStorage)
}

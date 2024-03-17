package agent

import (
	"io"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

type ClientMock struct {}

func (c ClientMock) Post (string, string, io.Reader) (*http.Response, error) {
    return &http.Response{}, nil
}


func TestSendMetrics(t *testing.T) {
	counterStorage := Storage[int64]{"SomeCounter": 1}
	gaugeStorage := Storage[float64]{"SomeGauge": 1.0}
	var client Client = ClientMock{}
	baseURL := "http://0.0.0.0:8080"

	SendMetrics(&counterStorage, &gaugeStorage, client, baseURL)()

	assert.Empty(t, counterStorage)
	assert.Empty(t, gaugeStorage)
}

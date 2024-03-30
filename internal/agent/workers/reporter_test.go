package workers

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/internal/types"
)

type ClientMock struct {}

func (c ClientMock) Post (string, string, io.Reader) (*http.Response, error) {
    return &http.Response{}, nil
}


func TestSendMetrics(t *testing.T) {
	counterStorage :=  map[string]types.Counter{"SomeCounter": 1}
	gaugeStorage :=  map[string]types.Gauge{"SomeGauge": 1.0}
	var client Client = ClientMock{}
	baseURL := "http://0.0.0.0:8080"

	SendMetrics(&counterStorage, &gaugeStorage, client, baseURL)()

	assert.Empty(t, counterStorage)
	assert.Empty(t, gaugeStorage)
}

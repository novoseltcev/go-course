package workers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	json "github.com/mailru/easyjson"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
)


type Client interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

func SendMetrics(counterStorage * map[string]model.Counter, gaugeStorage * map[string]model.Gauge, client Client, baseURL string) func() {
	fmt.Println("init SendMetrics worker")
	return func ()  {
		fmt.Printf("counters length=%d; gauge length=%d\n", len(*counterStorage), len(*gaugeStorage))

		for k, v := range *gaugeStorage {
			value := float64(v)
			err := send(client, baseURL, schema.Metrics{ID: k, MType: "gauge", Value: &value})
			if err == nil {
				delete(*gaugeStorage, k)
			}
		}

		for k, v := range *counterStorage {
			delta := int64(v)
			err := send(client, baseURL, schema.Metrics{ID: k, MType: "counter", Delta: &delta})
			if err == nil {
				delete(*counterStorage, k)
			}
		}
		fmt.Println("All sended")
	}
}

func send(client Client, baseURL string, metric schema.Metrics) error {
	var buf bytes.Buffer
	if _, err := json.MarshalToWriter(metric, &buf); err != nil {
		return err
	}
	
	url := baseURL + "/update/"
	response, err := client.Post(url, "application/json", &buf)
	if err != nil {
		fmt.Printf("Error during send request to %s\n", url)
		return err
	}
	
	if response.Body != nil {
		defer response.Body.Close()
		if _, err := io.Copy(io.Discard, response.Body); err != nil {
			fmt.Printf("Error during read request body from response to %s\n", url)
			return err
		}
	}

	fmt.Printf("Sended request to %s with code %s\n", url, response.Status)
	return nil
}

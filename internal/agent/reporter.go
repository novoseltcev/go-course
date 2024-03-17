package agent

import (
	"fmt"
	"io"
	"net/http"
)


type Client interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

func SendMetrics(counterStorage *Storage[int64], gaugeStorage *Storage[float64], client Client, baseUrl string) func() {
	fmt.Println("init SendMetrics worker")
	return func ()  {
		fmt.Printf("counters length=%d; gauge length=%d\n", len(*counterStorage), len(*gaugeStorage))

		for name, value := range *gaugeStorage {
			err := send(client, baseUrl, "gauge", name, fmt.Sprintf("%f", value))
			if err == nil {
				delete(*gaugeStorage, name)
			}
		}

		for name, value := range *counterStorage {
			err := send(client, baseUrl, "counter", name, fmt.Sprintf("%d", value))
			if err == nil {
				delete(*counterStorage, name)
			}
		}
		fmt.Println("All sended")
	}
}

func send(client Client, baseUrl string, metricType string, name string, value string) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", baseUrl, metricType, name, value)

	response, err := client.Post(url, "text/plain", http.NoBody)
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

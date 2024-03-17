package agent

import (
	"fmt"
	"net/http"
)


func SendMetrics(counterStorage *Storage[int64], gaugeStorage *Storage[float64], client *http.Client, baseUrl string) func() {
	fmt.Println("init SendMetrics worker")
	return func ()  {
		fmt.Printf("counters length=%d; gauge length=%d\n", len(*counterStorage), len(*gaugeStorage))

		for name, value := range *gaugeStorage {
			send(client, baseUrl, "gauge", name, fmt.Sprintf("%f", value))
		}

		for name, value := range *counterStorage {
			send(client, baseUrl, "counter", name, fmt.Sprintf("%d", value))
		}
		fmt.Println("All sended")
	}
}

func send(client *http.Client, baseUrl string, metricType string, name string, value string) {
	url := fmt.Sprintf("%s/update/%s/%s/%s", baseUrl, metricType, name, value)
	response, err := client.Post(url, "text/plain", http.NoBody)
	
	if err == nil {
		fmt.Printf("Sended request to %s with code %s\n", url, response.Status)
		return
	}

	fmt.Printf("Sended request to %s with code %s\n", url, response.Status)
}

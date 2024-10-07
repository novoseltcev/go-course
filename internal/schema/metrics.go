package schema

const (
	Gauge   = "gauge"
	Counter = "counter"
)

//go:generate easyjson -all metrics.go
type Metric struct {
	ID    string   `db:"id"    json:"id"`              // имя метрики
	MType string   `db:"type"  json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `db:"delta" json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `db:"value" json:"value,omitempty"` // значение метрики в случае передачи gauge
}

//easyjson:json
type MetricSlice []Metric

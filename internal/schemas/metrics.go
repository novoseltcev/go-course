// Package schemas contains seralizable and deserializable data structures.
package schemas

const (
	Gauge   = "gauge"
	Counter = "counter"
)

//go:generate easyjson -all metrics.go
type Metric struct {
	ID    string   `db:"id"    json:"id"`              // name of metric
	MType string   `db:"type"  json:"type"`            // parameter, specifying gauge or counter
	Delta *int64   `db:"delta" json:"delta,omitempty"` // value of metric in case of counter
	Value *float64 `db:"value" json:"value,omitempty"` // value of metric in case of gauge
}

//easyjson:json
type MetricSlice []Metric

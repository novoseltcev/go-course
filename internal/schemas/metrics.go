// Package schemas contains seralizable and deserializable data structures.
package schemas

import "errors"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

//go:generate easyjson -all metrics.go
type Metric struct {
	ID    string   `db:"name"  json:"id"`              // name of metric
	MType string   `db:"type"  json:"type"`            // parameter, specifying gauge or counter
	Delta *int64   `db:"delta" json:"delta,omitempty"` // value of metric in case of counter
	Value *float64 `db:"value" json:"value,omitempty"` // value of metric in case of gauge
}

// nolint: err113
func (m *Metric) Validate() error {
	if m.ID == "" {
		return errors.New("id is required")
	}

	if m.MType != Gauge && m.MType != Counter {
		return errors.New("type is invalid")
	}

	if m.MType == Gauge {
		if m.Value == nil {
			return errors.New("gauge metric has nil value")
		}

		if m.Delta != nil {
			return errors.New("gauge metric has non-nil delta")
		}
	}

	if m.MType == Counter {
		if m.Delta == nil {
			return errors.New("counter metric has nil delta")
		}

		if m.Value != nil {
			return errors.New("counter metric has non-nil value")
		}
	}

	return nil
}

//easyjson:json
type MetricSlice []Metric

//easyjson:json
type MetricIdentifier struct {
	ID    string `db:"id"   json:"id"`   // name of metric
	MType string `db:"type" json:"type"` // parameter, specifying gauge or counter
}

//nolint: err113
func (m *MetricIdentifier) Validate() error {
	if m.ID == "" {
		return errors.New("id is empty")
	}

	if m.MType != Gauge && m.MType != Counter {
		return errors.New("type is invalid")
	}

	return nil
}

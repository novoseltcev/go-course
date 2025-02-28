// Package schemas contains seralizable and deserializable data structures.
package schemas

import (
	"errors"
	"fmt"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var (
	ErrEmptyID      = errors.New("id is empty")
	ErrInvalidType  = errors.New("type is invalid")
	ErrInvalidValue = errors.New("value is invalid")
	ErrInvalidDelta = errors.New("delta is invalid")
)

//go:generate easyjson -all metrics.go

//easyjson:json
type MetricIdentifier struct {
	ID    string `db:"id"   json:"id"`   // name of metric
	MType string `db:"type" json:"type"` // parameter, specifying gauge or counter
}

func (m *MetricIdentifier) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("metric validator: %w", ErrEmptyID)
	}

	if m.MType != Gauge && m.MType != Counter {
		return fmt.Errorf("metric validator: %w", ErrInvalidType)
	}

	return nil
}

//easyjson:json
type MetricValues struct {
	Value *float64 `db:"value" json:"value,omitempty"` // value of metric in case of gauge
	Delta *int64   `db:"delta" json:"delta,omitempty"` // value of metric in case of counter
}

//easyjson:json
type Metric struct {
	Delta *int64   `db:"delta" json:"delta,omitempty"` // value of metric in case of counter
	Value *float64 `db:"value" json:"value,omitempty"` // value of metric in case of gauge
	ID    string   `db:"name"  json:"id"`              // name of metric
	MType string   `db:"type"  json:"type"`            // parameter, specifying gauge or counter
}

func (m *Metric) Validate() error {
	mID := MetricIdentifier{ID: m.ID, MType: m.MType}
	if err := mID.Validate(); err != nil {
		return err
	}

	if m.MType == Gauge {
		if m.Value == nil {
			return fmt.Errorf("metric validator: %w; want non-nil", ErrInvalidValue)
		}

		if m.Delta != nil {
			return fmt.Errorf("metric validator: %w; want nil", ErrInvalidDelta)
		}
	}

	if m.MType == Counter {
		if m.Delta == nil {
			return fmt.Errorf("metric validator: %w; want non-nil", ErrInvalidDelta)
		}

		if m.Value != nil {
			return fmt.Errorf("metric validator: %w; want nil", ErrInvalidValue)
		}
	}

	return nil
}

//easyjson:json
type MetricSlice []Metric

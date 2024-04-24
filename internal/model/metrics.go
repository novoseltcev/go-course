package model

type Metric struct {
	Name string
	Type string
	Value *float64
	Delta *int64
}

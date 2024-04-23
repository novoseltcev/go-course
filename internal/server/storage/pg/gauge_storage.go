package pg

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/model"
)


type GaugeStorage struct {
	DB *sqlx.DB
}

func (s *GaugeStorage) GetByName(ctx context.Context, name string) *model.Gauge {
	row := s.DB.QueryRowxContext(ctx, "SELECT name, value FROM metrics WHERE type = 'gauge' AND name = $1", name)
	if row == nil {
		return nil
	}
	var result model.Gauge
	if err := row.Scan(&result); err != nil {
		panic(err)
	}
	return &result
}

func (s *GaugeStorage) GetAll(ctx context.Context) []model.Metric[model.Gauge] {
	rows, err := s.DB.QueryxContext(ctx, "SELECT name, value FROM metrics WHERE type = 'gauge'")
	if err != nil {
		panic(err)
	}

	var metrics []model.Metric[model.Gauge]
	for rows.Next() {
		var metric model.Metric[model.Gauge]
		if err := rows.StructScan(&metric); err != nil {
			panic(err)
		}
		metrics = append(metrics, metric)
	}
	
	return metrics
}

func (s *GaugeStorage) Update(ctx context.Context, name string, value model.Gauge) {
	s.DB.MustExecContext(
		ctx,
		`INSERT INTO metrics ("name", "value", "type") VALUES ($1, $2, 'gauge') ON CONFLICT (name, type) DO UPDATE SET value = EXCLUDED.value`,
		name,
		value,
	)
}

package pg

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/model"
)


type CounterStorage struct {
	DB *sqlx.DB
}

func (s *CounterStorage) GetByName(ctx context.Context, name string) *model.Counter {
	row := s.DB.QueryRowxContext(ctx, "SELECT value FROM metrics WHERE type = 'counter' AND name = $1", name)
	if row == nil {
		return nil
	}
	var result model.Counter
	if err := row.Scan(&result); err != nil {
		panic(err)
	}
	return &result
}

func (s *CounterStorage) GetAll(ctx context.Context) []model.Metric[model.Counter] {
	rows, err := s.DB.QueryxContext(ctx, "SELECT name, value::integer FROM metrics")
	if err != nil {
		panic(err)
	}

	var metrics []model.Metric[model.Counter]
	for rows.Next() {
		var metric model.Metric[model.Counter]
		if err := rows.StructScan(&metric); err != nil {
			panic(err)
		}
		metrics = append(metrics, metric)
	}
	
	return metrics
}

func (s *CounterStorage) Update(ctx context.Context, name string, value model.Counter) {
	s.DB.MustExecContext(
		ctx,
		`INSERT INTO metrics ("name", "value", "type") VALUES ($1, $2, 'counter') ON CONFLICT (name, type) DO UPDATE SET value = EXCLUDED.value`,
		name,
		value,
	)
}

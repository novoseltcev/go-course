package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/model"
)


type Storage struct {
	DB *sqlx.DB
}

func (s *Storage) GetByName(ctx context.Context, name, Type string) (*model.Metric, error) {
	var result model.Metric
	err := s.DB.Get(&result, "SELECT name, type, value, delta FROM metrics WHERE type = $1 AND name = $2", Type, name)
	if errors.Is(err, sql.ErrNoRows) {
		return &result, nil
	}
	return &result, err
}

func (s *Storage) GetAll(ctx context.Context) ([]model.Metric, error) {
	rows, err := s.DB.QueryxContext(ctx, "SELECT name, type, value, delta FROM metrics")
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var metrics []model.Metric
	for rows.Next() {
		var metric model.Metric
		if err := rows.StructScan(&metric); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	
	return metrics, nil
}

func (s *Storage) Save(ctx context.Context, metric model.Metric) error {
	stmt, err := s.DB.PrepareContext(ctx, `INSERT INTO metrics (name, type, value, delta) VALUES ($1, $2, $3, $4) ON CONFLICT (name, type) DO UPDATE SET value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta`)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, metric.Name, metric.Type, metric.Value, metric.Delta); err != nil {
		return err
	}
	return nil
}

func (s *Storage) SaveAll(ctx context.Context, metrics []model.Metric) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO metrics (name, type, value, delta) VALUES ($1, $2, $3, $4) ON CONFLICT (name, type) DO UPDATE SET value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta`)
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		_, err := stmt.ExecContext(ctx, metric.Name, metric.Type, metric.Value, metric.Delta)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(err, rollbackErr)
			}
			return err
		}
	}

	return tx.Commit()
}

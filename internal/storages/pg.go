package storages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/schemas"
)

type PgStorage struct {
	DB *sqlx.DB
}

func NewPgStorage(url string) (*PgStorage, error) {
	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	return &PgStorage{DB: db}, nil
}

func (s *PgStorage) GetByName(ctx context.Context, name, metricType string) (*schemas.Metric, error) {
	var result schemas.Metric
	err := s.DB.GetContext(
		ctx,
		&result,
		"SELECT name, type, value, delta FROM metrics WHERE type = $1 AND name = $2",
		metricType,
		name,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &result, nil
	}

	return &result, err
}

func (s *PgStorage) GetAll(ctx context.Context) ([]schemas.Metric, error) {
	var metrics []schemas.Metric
	err := s.DB.SelectContext(ctx, &metrics, "SELECT name, type, value, delta FROM metrics")

	return metrics, err
}

func (s *PgStorage) Save(ctx context.Context, metric *schemas.Metric) error {
	stmt, err := s.DB.PrepareContext(
		ctx,
		`INSERT INTO metrics (name, type, value, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name, type) DO UPDATE SET value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta`,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Value, metric.Delta); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) SaveAll(ctx context.Context, metrics []schemas.Metric) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO metrics (name, type, value, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name, type) DO UPDATE SET value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, metric := range metrics {
		_, err := stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Value, metric.Delta)
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}

	return tx.Commit()
}

func (s *PgStorage) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *PgStorage) Close() error {
	return s.DB.Close()
}

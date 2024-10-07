package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/schema"
	s "github.com/novoseltcev/go-course/internal/server/storage"
)

type storage struct {
	DB *sqlx.DB
}

func New(url string) (s.MetricStorager, error) {
	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	return &storage{DB: db}, nil
}

func (s *storage) GetByName(ctx context.Context, name, metricType string) (*schema.Metric, error) {
	var result schema.Metric
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

func (s *storage) GetAll(ctx context.Context) ([]schema.Metric, error) {
	var metrics []schema.Metric
	err := s.DB.SelectContext(ctx, &metrics, "SELECT name, type, value, delta FROM metrics")

	return metrics, err
}

func (s *storage) Save(ctx context.Context, metric *schema.Metric) error {
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

func (s *storage) SaveAll(ctx context.Context, metrics []schema.Metric) error {
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

func (s *storage) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *storage) Close() error {
	return s.DB.Close()
}

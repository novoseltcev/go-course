package storages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/schemas"
)

type PgStorage struct {
	db *sqlx.DB
}

func NewPgStorage(db *sqlx.DB) *PgStorage {
	return &PgStorage{db: db}
}

func (s *PgStorage) GetOne(ctx context.Context, id, mType string) (*schemas.Metric, error) {
	var result schemas.Metric
	err := s.db.GetContext(
		ctx,
		&result,
		"SELECT name, type, value, delta FROM metrics WHERE type = $1 AND name = $2",
		mType,
		id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Join(err, ErrNotFound)
	}

	return &result, err
}

func (s *PgStorage) GetAll(ctx context.Context) ([]schemas.Metric, error) {
	var metrics []schemas.Metric
	err := s.db.SelectContext(ctx, &metrics, "SELECT name, type, value, delta FROM metrics")

	return metrics, err
}

func (s *PgStorage) Save(ctx context.Context, metric *schemas.Metric) error {
	stmt, err := s.db.PrepareContext(
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

func (s *PgStorage) SaveBatch(ctx context.Context, metrics []schemas.Metric) error {
	tx, err := s.db.BeginTx(ctx, nil)
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

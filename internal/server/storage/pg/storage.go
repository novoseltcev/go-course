package pg

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/model"
	s "github.com/novoseltcev/go-course/internal/server/storage"
)


type storage struct {
	DB *sqlx.DB
}



func New(URL string) (s.MetricStorager, error) {
	db, err := sqlx.Open("pgx", URL)
	if err != nil {
		return nil, err
	}

	return &storage{DB: db}, nil
}

func (s *storage) GetByName(ctx context.Context, name, Type string) (*model.Metric, error) {
	var result model.Metric
	err := s.DB.GetContext(ctx, &result, "SELECT name, type, value, delta FROM metrics WHERE type = $1 AND name = $2", Type, name)
	if errors.Is(err, sql.ErrNoRows) {
		return &result, nil
	}
	return &result, err
}

func (s *storage) GetAll(ctx context.Context) ([]model.Metric, error) {
	var metrics []model.Metric
	err := s.DB.SelectContext(ctx, &metrics, "SELECT name, type, value, delta FROM metrics")
	return metrics, err
}

func (s *storage) Save(ctx context.Context, metric model.Metric) error {
	stmt, err := s.DB.PrepareContext(ctx, `INSERT INTO metrics (name, type, value, delta) VALUES ($1, $2, $3, $4) ON CONFLICT (name, type) DO UPDATE SET value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta`)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, metric.Name, metric.Type, metric.Value, metric.Delta); err != nil {
		return err
	}
	return nil
}

func (s *storage) SaveAll(ctx context.Context, metrics []model.Metric) error {
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

func (s *storage) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}


func (s *storage) Close() error {
	return s.DB.Close()
}

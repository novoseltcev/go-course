package storages_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

var (
	migrationsDir = filepath.Join("..", "..", "migrations")
	ctxTimeout    = 10 * time.Second // nolint:mnd
	testID        = "test"
)

func getMetrics(ctx context.Context, t *testing.T, db *sqlx.DB) []schemas.Metric {
	t.Helper()

	var result []schemas.Metric
	require.NoError(t, db.SelectContext(ctx, &result, "SELECT name, type, value, delta FROM metrics"))

	return result
}

// nolint: paralleltest  // idk why differrent containers influences on each other
func TestPgStorage_GetOne_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	t.Cleanup(cancel)

	storage := storages.NewPgStorage(helpers.SetupDB(ctx, t, migrationsDir, "base.sql"))

	tests := []struct {
		name string
		got  schemas.MetricIdentifier
		want schemas.Metric
	}{
		{
			name: "gauge",
			got: schemas.MetricIdentifier{
				ID:    testID,
				MType: schemas.Gauge,
			},
			want: schemas.Metric{
				ID:    testID,
				MType: schemas.Gauge,
				Value: &value,
			},
		},
		{
			name: "counter",
			got: schemas.MetricIdentifier{
				ID:    testID,
				MType: schemas.Counter,
			},
			want: schemas.Metric{
				ID:    testID,
				MType: schemas.Counter,
				Delta: &delta,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.GetOne(ctx, tt.got.ID, tt.got.MType)
			require.NoError(t, err)
			assert.Equal(t, &tt.want, got)
		})
	}
}

func TestPgStorage_GetOne_NotFound(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	t.Cleanup(cancel)

	storage := storages.NewPgStorage(helpers.SetupDB(ctx, t, migrationsDir, "base.sql"))

	tests := []struct {
		name string
		got  schemas.MetricIdentifier
	}{
		{
			name: "unknown gauge",
			got: schemas.MetricIdentifier{
				ID:    testutils.UNKNOWN,
				MType: schemas.Gauge,
			},
		},
		{
			name: "unknown counter",
			got: schemas.MetricIdentifier{
				ID:    testutils.UNKNOWN,
				MType: schemas.Counter,
			},
		},
		{
			name: "unknown type",
			got: schemas.MetricIdentifier{
				ID:    testID,
				MType: testutils.UNKNOWN,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := storage.GetOne(ctx, tt.got.ID, tt.got.MType)
			assert.ErrorIs(t, err, storages.ErrNotFound)
		})
	}
}

func TestPgStorage_GetAll_Empty(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	t.Cleanup(cancel)

	storage := storages.NewPgStorage(helpers.SetupDB(ctx, t, migrationsDir))

	metrics, err := storage.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, metrics)
}

// nolint: paralleltest  // idk why differrent containers influences on each other
func TestPgStorage_GetAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	t.Cleanup(cancel)

	storage := storages.NewPgStorage(helpers.SetupDB(ctx, t, migrationsDir, "base.sql"))

	metrics, err := storage.GetAll(ctx)
	require.NoError(t, err)
	assert.ElementsMatch(t, []schemas.Metric{
		{
			ID:    testID,
			MType: schemas.Gauge,
			Value: &value,
		},
		{
			ID:    testID,
			MType: schemas.Counter,
			Delta: &delta,
		},
	}, metrics)
}

func TestPgStorage_Save(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	t.Cleanup(cancel)

	db := helpers.SetupDB(ctx, t, migrationsDir)
	storage := storages.NewPgStorage(db)

	got := schemas.Metric{ID: testutils.STRING, MType: schemas.Gauge, Value: &value}
	require.NoError(t, storage.Save(ctx, &got))
	assert.Equal(t, []schemas.Metric{got}, getMetrics(ctx, t, db))
}

func TestPgStorage_Save_Exists(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	t.Cleanup(cancel)

	db := helpers.SetupDB(ctx, t, migrationsDir, "base.sql")
	storage := storages.NewPgStorage(db)

	t.Run("update gauge", func(t *testing.T) {
		t.Parallel()

		got := schemas.Metric{ID: testID, MType: schemas.Gauge, Value: &newValue}

		require.NoError(t, storage.Save(ctx, &got))
		assert.Contains(t, getMetrics(ctx, t, db), got)
	})

	t.Run("update counter", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, storage.Save(ctx, &schemas.Metric{
			ID:    testID,
			MType: schemas.Counter,
			Delta: &inc,
		}))

		assert.Contains(t, getMetrics(ctx, t, db), schemas.Metric{
			ID:    testID,
			MType: schemas.Counter,
			Delta: &newDelta,
		})
	})
}

func TestPgStorage_SaveBatch(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	t.Cleanup(cancel)

	db := helpers.SetupDB(ctx, t, migrationsDir, "base.sql")
	storage := storages.NewPgStorage(db)

	require.NoError(t,
		storage.SaveBatch(ctx, []schemas.Metric{
			// add
			{
				ID:    testutils.STRING,
				MType: schemas.Gauge,
				Value: &value,
			},
			{
				ID:    testutils.STRING,
				MType: schemas.Counter,
				Delta: &delta,
			},
			// update
			{
				ID:    testID,
				MType: schemas.Gauge,
				Value: &newValue,
			},
			{
				ID:    testID,
				MType: schemas.Counter,
				Delta: &inc,
			},
		}),
	)

	assert.ElementsMatch(t, []schemas.Metric{
		{
			ID:    testutils.STRING,
			MType: schemas.Gauge,
			Value: &value,
		},
		{
			ID:    testutils.STRING,
			MType: schemas.Counter,
			Delta: &delta,
		},
		{
			ID:    testID,
			MType: schemas.Gauge,
			Value: &newValue,
		},
		{
			ID:    testID,
			MType: schemas.Counter,
			Delta: &newDelta,
		},
	}, getMetrics(ctx, t, db))
}

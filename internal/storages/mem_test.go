package storages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

var (
	value    = 10.123
	newValue = float64(5)
	delta    = int64(10)
	inc      = int64(2)
	newDelta = delta + inc
)

func setupMemStorage(t *testing.T) *storages.MemStorage {
	t.Helper()

	storage := storages.NewMemStorage()
	storage.Data[schemas.Gauge][testutils.STRING] = schemas.MetricValues{Value: &value}
	storage.Data[schemas.Counter][testutils.STRING] = schemas.MetricValues{Delta: &delta}

	return storage
}

func TestMemStorage_GetOne(t *testing.T) {
	t.Parallel()

	storage := setupMemStorage(t)

	tests := []struct {
		name string
		got  schemas.MetricIdentifier
		want schemas.Metric
	}{
		{
			name: "gauge",
			got: schemas.MetricIdentifier{
				ID:    testutils.STRING,
				MType: schemas.Gauge,
			},
			want: schemas.Metric{
				ID:    testutils.STRING,
				MType: schemas.Gauge,
				Value: &value,
			},
		},
		{
			name: "counter",
			got: schemas.MetricIdentifier{
				ID:    testutils.STRING,
				MType: schemas.Counter,
			},
			want: schemas.Metric{
				ID:    testutils.STRING,
				MType: schemas.Counter,
				Delta: &delta,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := storage.GetOne(context.TODO(), tt.got.ID, tt.got.MType)

			require.NoError(t, err)
			assert.Equal(t, &tt.want, got)
		})
	}
}

func TestMemStorage_GetOne_NotFound(t *testing.T) {
	t.Parallel()

	storage := setupMemStorage(t)

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
				ID:    testutils.STRING,
				MType: testutils.UNKNOWN,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := storage.GetOne(context.TODO(), tt.got.ID, tt.got.MType)

			assert.ErrorIs(t, err, storages.ErrNotFound)
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	t.Parallel()

	storage := setupMemStorage(t)

	got, err := storage.GetAll(context.TODO())
	require.NoError(t, err)
	assert.Len(t, got, 2)

	assert.ElementsMatch(t, got, []schemas.Metric{
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
	})
}

func TestMemStorage_Save(t *testing.T) {
	t.Parallel()

	storage := storages.NewMemStorage()

	require.NoError(t,
		storage.Save(context.TODO(), &schemas.Metric{
			ID:    testutils.STRING,
			MType: schemas.Gauge,
			Value: &value,
		}),
	)

	assert.Empty(t, storage.Data[schemas.Counter])
	assert.Len(t, storage.Data[schemas.Gauge], 1)
	assert.Equal(t, schemas.MetricValues{Value: &value}, storage.Data[schemas.Gauge][testutils.STRING])
}

func TestMemStorage_Save_Exists(t *testing.T) {
	t.Parallel()

	storage := setupMemStorage(t)

	t.Run("update gauge", func(t *testing.T) {
		t.Parallel()

		require.NoError(t,
			storage.Save(context.TODO(), &schemas.Metric{
				ID:    testutils.STRING,
				MType: schemas.Gauge,
				Value: &newValue,
			}),
		)

		assert.Len(t, storage.Data[schemas.Gauge], 1)
		assert.Equal(t, schemas.MetricValues{Value: &newValue}, storage.Data[schemas.Gauge][testutils.STRING])
	})

	t.Run("update counter", func(t *testing.T) {
		t.Parallel()

		require.NoError(t,
			storage.Save(context.TODO(), &schemas.Metric{
				ID:    testutils.STRING,
				MType: schemas.Counter,
				Delta: &inc,
			}),
		)

		assert.Len(t, storage.Data[schemas.Counter], 1)
		assert.Equal(t, schemas.MetricValues{Delta: &newDelta}, storage.Data[schemas.Counter][testutils.STRING])
	})
}

func TestMemStorage_SaveBatch(t *testing.T) {
	t.Parallel()

	storage := storages.NewMemStorage()

	require.NoError(t,
		storage.SaveBatch(context.TODO(), []schemas.Metric{
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
		}),
	)

	assert.Len(t, storage.Data[schemas.Counter], 1)
	assert.Len(t, storage.Data[schemas.Gauge], 1)
	assert.Equal(t, schemas.MetricValues{Value: &value}, storage.Data[schemas.Gauge][testutils.STRING])
	assert.Equal(t, schemas.MetricValues{Delta: &delta}, storage.Data[schemas.Counter][testutils.STRING])
}

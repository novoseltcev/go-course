package collectors_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/agent/collectors"
	"github.com/novoseltcev/go-course/internal/schemas"
)

func TestCollectCoreMetrics(t *testing.T) {
	t.Parallel()

	res, err := collectors.CollectCoreMetrics(context.TODO())

	require.NoError(t, err)
	assert.Len(t, res, 3)

	keys := make([]string, 0, 3)

	for _, m := range res {
		keys = append(keys, m.ID)
		assert.Equal(t, schemas.Gauge, m.MType)
		assert.NotNil(t, m.Value)
		assert.Nil(t, m.Delta)
	}

	assert.EqualValues(t, []string{"TotalMemory", "FreeMemory", "CPUutilization1"}, keys)
}

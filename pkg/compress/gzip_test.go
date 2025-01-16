package compress_test

import (
	"compress/gzip"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/compress"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	for i := gzip.HuffmanOnly; i < gzip.BestCompression; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			cmp, err := compress.NewGzip(i)
			require.NoError(t, err)

			data := []byte("test")

			compressed, err := cmp.Compress(data)
			require.NoError(t, err)

			decompressed, err := cmp.Decompress(compressed)
			require.NoError(t, err)

			assert.NotEqual(t, data, compressed)
			assert.Equal(t, data, decompressed)
		})
	}
}

func TestNewGzipLevelError(t *testing.T) {
	t.Parallel()

	t.Run("-3", func(t *testing.T) {
		t.Parallel()

		_, err := compress.NewGzip(-3)
		assert.EqualError(t, err, "invalid compression level: -3")
	})

	t.Run("10", func(t *testing.T) {
		t.Parallel()

		_, err := compress.NewGzip(10)
		assert.EqualError(t, err, "invalid compression level: 10")
	})
}

package compress_test

import (
	"compress/gzip"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/compress"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

func TestSuccessCompressAndDecompress(t *testing.T) {
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

func TestCompress_ErrorInvalidLevel(t *testing.T) {
	t.Parallel()

	for _, level := range []int{-3, 10} {
		t.Run(strconv.Itoa(level), func(t *testing.T) {
			t.Parallel()

			cmp := compress.GzipCompressor{Level: level}
			_, err := cmp.Compress(testutils.Bytes)

			assert.EqualError(t, err, fmt.Sprintf("gzip: invalid compression level: %d", level))
		})
	}
}

func TestDecompess_ErrorCreateReader(t *testing.T) {
	t.Parallel()

	cmp, err := compress.NewGzip(gzip.BestCompression)
	require.NoError(t, err)

	_, err = cmp.Decompress([]byte{})

	assert.ErrorIs(t, err, io.EOF)
}

func TestDecompess_ErrorCopy(t *testing.T) {
	t.Parallel()

	cmp, err := compress.NewGzip(gzip.BestCompression)
	require.NoError(t, err)

	_, err = cmp.Decompress([]byte{})

	assert.ErrorIs(t, err, io.EOF)
}

func TestNewGzip_ErrorInvalidLevel(t *testing.T) {
	t.Parallel()

	for _, level := range []int{-3, 10} {
		t.Run(strconv.Itoa(level), func(t *testing.T) {
			t.Parallel()

			_, err := compress.NewGzip(level)
			assert.EqualError(t, err, fmt.Sprintf("invalid compression level: %d", level))
		})
	}
}

package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa/utils"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

func TestSplitToChunks(t *testing.T) {
	t.Parallel()

	data := []int{1, 2, 3, 4, 5}
	chunks, err := utils.SplitToChunks(data, 2)
	require.NoError(t, err)

	require.Len(t, chunks, 3)
	assert.Equal(t, []int{1, 2}, chunks[0])
	assert.Equal(t, []int{3, 4}, chunks[1])
	assert.Equal(t, []int{5}, chunks[2])
}

func TestSplitToChunksWithRemainder(t *testing.T) {
	t.Parallel()

	data := []int{1, 2, 3, 4, 5, 6}
	chunks, err := utils.SplitToChunks(data, 2)
	require.NoError(t, err)

	require.Len(t, chunks, 3)
	assert.Equal(t, []int{1, 2}, chunks[0])
	assert.Equal(t, []int{3, 4}, chunks[1])
	assert.Equal(t, []int{5, 6}, chunks[2])
}

func TestSplitToChunksWithZeroChunkSizeFailsErr(t *testing.T) {
	t.Parallel()

	_, err := utils.SplitToChunks(testutils.Bytes, 0)
	require.ErrorIs(t, utils.ErrInvalidChunkSize, err)
}

func ExampleSplitToChunks() {
	chunks, _ := utils.SplitToChunks([]int{1, 2, 3, 4, 5}, 2)
	fmt.Println(chunks)

	// Output:
	// [[1 2] [3 4] [5]]
}

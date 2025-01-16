package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa/utils"
)

func TestSplitToChunks(t *testing.T) {
	t.Parallel()

	data := []int{1, 2, 3, 4, 5}
	chunks := utils.SplitToChunks(data, 2)

	require.Len(t, chunks, 3)
	assert.Equal(t, []int{1, 2}, chunks[0])
	assert.Equal(t, []int{3, 4}, chunks[1])
	assert.Equal(t, []int{5}, chunks[2])
}

func TestSplitToChunksWithRemainder(t *testing.T) {
	t.Parallel()

	data := []int{1, 2, 3, 4, 5, 6}
	chunks := utils.SplitToChunks(data, 2)

	require.Len(t, chunks, 3)
	assert.Equal(t, []int{1, 2}, chunks[0])
	assert.Equal(t, []int{3, 4}, chunks[1])
	assert.Equal(t, []int{5, 6}, chunks[2])
}

func ExampleSplitToChunks() {
	chunks := utils.SplitToChunks([]int{1, 2, 3, 4, 5}, 2)
	fmt.Println(chunks)

	// Output:
	// [[1 2] [3 4] [5]]
}

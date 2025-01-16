package utils

// SplitToChunks splits data into chunks of the size of chunkSize.
func SplitToChunks[T any](data []T, chunkSize int) [][]T {
	chunks := make([][]T, 0)

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunks = append(chunks, data[i:end])
	}

	return chunks
}

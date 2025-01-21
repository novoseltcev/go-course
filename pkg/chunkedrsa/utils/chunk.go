package utils

import "errors"

var ErrInvalidChunkSize = errors.New("chunk size must be positive")

// SplitToChunks splits data into chunks of the size of chunkSize.
func SplitToChunks[T any](data []T, chunkSize int) ([][]T, error) {
	if chunkSize <= 0 {
		return nil, ErrInvalidChunkSize
	}

	chunks := make([][]T, 0)

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunks = append(chunks, data[i:end])
	}

	return chunks, nil
}

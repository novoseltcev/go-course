package compress

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

type GzipCompressor struct {
	Level int
}

var ErrCompressionLevel = errors.New("invalid compression level")

func NewGzip(level int) (*GzipCompressor, error) {
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		return nil, fmt.Errorf("%w: %d", ErrCompressionLevel, level)
	}

	return &GzipCompressor{Level: level}, nil
}

func (gc *GzipCompressor) Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	gzw, err := gzip.NewWriterLevel(buf, gc.Level)
	if err != nil {
		return nil, err
	}

	if _, err = gzw.Write(data); err != nil {
		return nil, err
	}

	if err = gzw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (gc *GzipCompressor) Decompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	gzr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	for {
		if _, err := io.CopyN(buf, gzr, 1024); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}
	}

	return buf.Bytes(), nil
}

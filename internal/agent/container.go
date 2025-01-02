package agent

import (
	"compress/gzip"
	"crypto/sha256"
	"net/http"
	"time"

	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/pkg/compress"
	"github.com/novoseltcev/go-course/pkg/cryptoalg"
	"github.com/novoseltcev/go-course/pkg/hash"
	"github.com/novoseltcev/go-course/pkg/retry"
)

type AppContainer struct {
	Cfg      *Config
	Reporter reporters.Reporter
}

func NewAppContainer(cfg *Config) (*AppContainer, error) {
	compressor, err := compress.NewGzip(gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	opts := []reporters.Option{
		reporters.WithCompression(compressor),
		reporters.WithRetry(retry.Options{
			Retries:  3, //nolint:mnd
			Attempts: []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
		}),
	}

	if cfg.CryptoKey != "" {
		enc, err := cryptoalg.NewPKCS1v15EncryptorFromFile(cfg.CryptoKey)
		if err != nil {
			return nil, err
		}

		opts = append(opts, reporters.WithEncryption(enc))
	}

	if cfg.SecretKey != "" {
		opts = append(opts, reporters.WithCheckSum(hash.NewHMAC(cfg.SecretKey, sha256.New)))
	}

	return &AppContainer{
		Cfg:      cfg,
		Reporter: reporters.NewHTTPClient(http.DefaultClient, cfg.Address, opts...),
	}, nil
}

func (c *AppContainer) Close() error {
	return nil
}

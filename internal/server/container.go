package server

import (
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/cryptoalg"
)

type AppContainer struct {
	Cfg       *Config
	Storage   storages.MetricStorager
	Decryptor *cryptoalg.PKCS1v15Algorithm
}

func NewAppContainer(cfg *Config) (*AppContainer, error) {
	var err error

	container := &AppContainer{Cfg: cfg}

	if cfg.DatabaseDsn == "" {
		container.Storage = storages.NewMemStorage()
	} else {
		container.Storage, err = storages.NewPgStorage(cfg.DatabaseDsn)
		if err != nil {
			return nil, err
		}
	}

	if cfg.CryptoKey != "" {
		container.Decryptor, err = cryptoalg.NewPKCS1v15DecryptorFromFile(cfg.CryptoKey)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Restore && cfg.FileStoragePath != "" {
		if err := restore(cfg.FileStoragePath, container.Storage); err != nil {
			return nil, err
		}
	}

	return container, nil
}

func (c *AppContainer) Close() error {
	return c.Storage.Close()
}

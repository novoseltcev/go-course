// nolint: tagliatelle
package server

import (
	"encoding/json"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

type Config struct {
	Address          string        `env:"ADDRESS"           json:"address"`
	Restore          bool          `env:"RESTORE"           json:"restore"`
	FileStoragePath  string        `env:"FILE_STORAGE_PATH" json:"store_file"`
	RawStoreInterval string        `env:"STORE_INTERVAL"    json:"store_interval"`
	DatabaseDsn      string        `env:"DATABASE_DSN"      json:"database_dsn"`
	SecretKey        string        `env:"KEY"               json:"-"`
	CryptoKey        string        `env:"CRYPTO_KEY,file"   json:"crypto_key"`
	StoreInterval    time.Duration `json:"-"`
}

func (c *Config) Load(fs afero.Fs, path string, flags *pflag.FlagSet) error {
	if path != "" {
		fd, err := fs.Open(path)
		if err != nil {
			return err
		}
		defer fd.Close()

		if err := json.NewDecoder(fd).Decode(c); err != nil {
			return err
		}
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}

	if err := env.Parse(c); err != nil {
		return err
	}

	if err := c.parseRawFields(); err != nil {
		return err
	}

	return nil
}

func (c *Config) parseRawFields() error {
	var err error
	if c.RawStoreInterval != "" {
		c.StoreInterval, err = time.ParseDuration(c.RawStoreInterval)
	}

	return err
}

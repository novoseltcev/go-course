//nolint: tagliatelle
package server

import "time"

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

func (c *Config) FinishParse() error {
	var err error
	c.StoreInterval, err = time.ParseDuration(c.RawStoreInterval)

	return err
}

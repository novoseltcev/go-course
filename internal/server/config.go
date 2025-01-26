// nolint: tagliatelle
package server

import (
	"encoding/json"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

type Config struct {
	Address           string        `env:"ADDRESS"           json:"address"`
	GRPCAddress       string        `env:"GRPC_ADDRESS"      json:"grpc_address"`
	Restore           bool          `env:"RESTORE"           json:"restore"`
	FileStoragePath   string        `env:"FILE_STORAGE_PATH" json:"store_file"`
	RawStoreInterval  string        `env:"STORE_INTERVAL"    json:"store_interval"`
	DatabaseDsn       string        `env:"DATABASE_DSN"      json:"database_dsn"`
	SecretKey         string        `env:"KEY"               json:"-"`
	CryptoKey         string        `env:"CRYPTO_KEY,file"   json:"crypto_key"`
	RawTrustedSubnets []string      `env:"TRUSTED_SUBNETS"   json:"trusted_subnets"`
	TrustedSubnets    []net.IPNet   `json:"-"`
	StoreInterval     time.Duration `json:"-"`
}

func (c *Config) Load(fs afero.Fs, path string, flags *pflag.FlagSet, args []string) error {
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

	if err := flags.Parse(args); err != nil {
		return err
	}

	if err := env.Parse(c); err != nil {
		return err
	}

	if err := c.ParseRawFields(); err != nil {
		return err
	}

	return nil
}

func (c *Config) ParseRawFields() error {
	var err error
	if c.RawStoreInterval != "" {
		c.StoreInterval, err = time.ParseDuration(c.RawStoreInterval)
		if err != nil {
			return err
		}
	}

	if c.RawTrustedSubnets != nil {
		c.TrustedSubnets = make([]net.IPNet, 0, len(c.RawTrustedSubnets))
		for _, rawSubnet := range c.RawTrustedSubnets {
			_, subnet, err := net.ParseCIDR(rawSubnet)
			if err != nil {
				return err
			}

			c.TrustedSubnets = append(c.TrustedSubnets, *subnet)
		}
	}

	return nil
}

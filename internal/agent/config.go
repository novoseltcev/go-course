// nolint: tagliatelle
package agent

import (
	"encoding/json"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

type Config struct {
	Address           string        `env:"ADDRESS"         json:"address"`
	RawPollInterval   string        `env:"POLL_INTERVAL"   json:"poll_interval"`
	RawReportInterval string        `env:"REPORT_INTERVAL" json:"report_interval"`
	SecretKey         string        `env:"KEY"             json:"-"`
	CryptoKey         string        `env:"CRYPTO_KEY,file" json:"crypto_key"`
	ReporterType      string        `env:"REPORTER_TYPE"   json:"reporter_type"`
	PollInterval      time.Duration `json:"-"`
	ReportInterval    time.Duration `json:"-"`
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
	if c.RawPollInterval != "" {
		c.PollInterval, err = time.ParseDuration(c.RawPollInterval)
		if err != nil {
			return err
		}
	}

	if c.RawReportInterval != "" {
		c.ReportInterval, err = time.ParseDuration(c.RawReportInterval)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint: tagliatelle
package agent

import "time"

type Config struct {
	Address           string        `env:"ADDRESS"         json:"address"`
	RawPollInterval   string        `env:"POLL_INTERVAL"   json:"poll_interval"`
	RawReportInterval string        `env:"REPORT_INTERVAL" json:"report_interval"`
	SecretKey         string        `env:"KEY"             json:"-"`
	CryptoKey         string        `env:"CRYPTO_KEY,file" json:"crypto_key"`
	PollInterval      time.Duration `json:"-"`
	ReportInterval    time.Duration `json:"-"`
}

func (c *Config) FinishParse() error {
	pollInterval, err := time.ParseDuration(c.RawPollInterval)
	if err != nil {
		return err
	}

	reportInterval, err := time.ParseDuration(c.RawReportInterval)
	if err != nil {
		return err
	}

	c.PollInterval = pollInterval
	c.ReportInterval = reportInterval

	return nil
}

package agent

import "time"

type Config struct {
	Address string
	PollInterval time.Duration
	ReportInterval time.Duration
}

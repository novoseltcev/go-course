package agent

import "time"

type Config struct {
	Address			string			`env:"ADDRESS"`
	PollInterval	time.Duration	`env:"POLL_INTERVAL"`
	RateLimit 		time.Duration 	`env:"REPORT_INTERVAL"`
	SecretKey		string			`env:"KEY"`
}

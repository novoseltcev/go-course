package retry

import "time"

type Options struct {
	Retries  uint
	Attempts []time.Duration
}

func (opt *Options) TotalAttempts() uint {
	if opt.Retries == 0 {
		return 3 //nolint:mnd
	}

	return opt.Retries
}

func (opt *Options) GetAttemptDelay(attempt uint) time.Duration {
	if len(opt.Attempts) == 0 {
		return time.Duration(attempt+1) * time.Second
	}

	if attempt >= uint(len(opt.Attempts)) {
		return opt.Attempts[len(opt.Attempts)-1]
	}

	return opt.Attempts[attempt]
}

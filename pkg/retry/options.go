package retry

import "time"

// Options is a configuration for retrying a function.
type Options struct {
	Attempts []time.Duration // list of timeouts for each attempt
	Retries  uint            // number of attempts
}

// TotalAttempts returns the total number of attempts.
//
// Default is 3.
func (opt *Options) TotalAttempts() uint {
	if opt.Retries == 0 {
		return 3 //nolint:mnd
	}

	return opt.Retries
}

// GetAttemptDelay returns the delay for the attempt.
//
// Default is attempt * 1 second.
// If the attempt is greater than the number of attempts, the last attempt is returned.
func (opt *Options) GetAttemptDelay(attempt uint) time.Duration {
	if len(opt.Attempts) == 0 {
		return time.Duration(attempt+1) * time.Second
	}

	if attempt >= uint(len(opt.Attempts)) {
		return opt.Attempts[len(opt.Attempts)-1]
	}

	return opt.Attempts[attempt]
}

package main

import (
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/novoseltcev/go-course/internal/agent"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Use this command to run agent",
		Run: func(cmd *cobra.Command, _ []string) {
			config, err := getConfig(cmd.Flags())
			if err != nil {
				log.Fatal(err)
			}

			config.PollInterval *= time.Second
			config.RateLimit *= time.Second

			a := agent.NewAgent(*config)
			ctx, _ := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			a.Start(ctx)
		},
	}

	flags := cmd.Flags()
	flags.StringP("a", "a", "localhost:8080", "Server address")
	flags.IntP("p", "p", 0, "poll runtime metrics interval")
	flags.IntP("r", "r", 0, "rate limit to send metrics")
	flags.StringP("k", "k", "", "Secret key for hashing data")

	return cmd
}

func getConfig(flags *pflag.FlagSet) (*agent.Config, error) {
	address, err := flags.GetString("a")
	if err != nil {
		return nil, err
	}

	pollInterval, err := flags.GetInt("p")
	if err != nil {
		return nil, err
	}

	rateLimit, err := flags.GetInt("r")
	if err != nil {
		return nil, err
	}

	secretKey, err := flags.GetString("k")
	if err != nil {
		return nil, err
	}

	config := agent.Config{
		Address:      address,
		PollInterval: time.Duration(pollInterval),
		RateLimit:    time.Duration(rateLimit),
		SecretKey:    secretKey,
	}

	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

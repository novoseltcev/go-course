package main

import (
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Use this command to run agent",
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			address, _ := flags.GetString("a")
			pollInterval, _ := flags.GetInt("p")
			rateLimit, _ := flags.GetInt("r")
			secretKey, _ := flags.GetString("k")

			config := agent.Config{
				Address: address,
				PollInterval: time.Duration(pollInterval),
				RateLimit: time.Duration(rateLimit),
				SecretKey: secretKey,
			}
			env.Parse(&config)
			config.PollInterval *= time.Second
			config.RateLimit *= time.Second

			a := agent.NewAgent(config)
			ctx, _ := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			a.Start(ctx)
		},
	}

	flags := cmd.Flags()
	flags.StringP("a", "a", "localhost:8080", "Server address")
	flags.IntP("p", "p", 2, "poll runtime metrics interval")
	flags.IntP("r", "r", 10, "rate limit to send metrics")
	flags.StringP("k", "k", "", "Secret key for hashing data")
	return cmd
}

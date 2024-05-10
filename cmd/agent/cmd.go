package main

import (
	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/spf13/cobra"
	"github.com/caarlos0/env/v10"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Use this command to run agent",
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			address, _ := flags.GetString("a")
			pollInterval, _ := flags.GetInt("p")
			reportInterval, _ := flags.GetInt("r")
			secretKey, _ := flags.GetString("k")

			config := agent.Config{
				Address: address,
				PollInterval: pollInterval,
				ReportInterval: reportInterval,
				SecretKey: secretKey,
			}
			env.Parse(&config)

			a := agent.NewAgent(config)
			a.Start()
		},
	}

	flags := cmd.Flags()
	flags.StringP("a", "a", "localhost:8080", "Server address")
	flags.IntP("p", "p", 2, "poll runtime metrics interval")
	flags.IntP("r", "r", 10, "send metrics to server interval")
	flags.StringP("k", "k", "", "Secret key for hashing data")
	return cmd
}

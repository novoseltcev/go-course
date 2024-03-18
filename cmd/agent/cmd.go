package main

import (
	"time"
	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Use this command to run agent",
		Run: func(cmd *cobra.Command, args []string) {
			address, _ := cmd.Flags().GetString("a")
			pollInterval, _ := cmd.Flags().GetInt("p")
			reportInterval, _ := cmd.Flags().GetInt("r")

			config := agent.Config{
				Address: address,
				PollInterval: time.Duration(pollInterval) * time.Second,
				ReportInterval: time.Duration(reportInterval) * time.Second,
			}
			a := agent.NewAgent(config)
			a.Start()
		},
	}

	flags := cmd.Flags()
	flags.StringP("a", "a", "localhost:8080", "Server address")
	flags.IntP("p", "p", 2, "poll runtime metrics interval")
	flags.IntP("r", "r", 10, "send metrics to server interval")
	return cmd
}

package main

import (
	"github.com/caarlos0/env/v10"
	"github.com/novoseltcev/go-course/internal/server"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Use this command to run server",
		Run: func(cmd *cobra.Command, args []string) {
			address, _ := cmd.Flags().GetString("a")
			config := server.Config{Address: address}
			env.Parse(&config)
	
			s := server.NewServer(config)
			if err := s.Start(); err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringP("a", "a", ":8080", "Server address")
	return cmd
}

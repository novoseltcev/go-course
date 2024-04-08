package main

import (
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/spf13/cobra"

	"github.com/novoseltcev/go-course/internal/server"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Use this command to run server",
		Run: func(cmd *cobra.Command, args []string) {
			address, _ := cmd.Flags().GetString("a")
			storeInterval, _ := cmd.Flags().GetInt("s")
			fileStoragePath, _ := cmd.Flags().GetString("f")
			restore, _ := cmd.Flags().GetBool("r")

			config := server.Config{
				Address: address,
				StoreInterval: time.Duration(storeInterval) * time.Second,
				FileStoragePath: fileStoragePath,
				Restore: restore,
			}
			env.Parse(&config)
	
			s := server.NewServer(config)
			if err := s.Start(); err != nil {
				panic(err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringP("a", "a", ":8080", "Server address")
	flags.IntP("s", "s", 300, "backup interval")
	flags.StringP("f", "f", "/tmp/metrics-db.json", "path to backup")
	flags.BoolP("r", "r", true, "restore from backup after restart")
	return cmd
}

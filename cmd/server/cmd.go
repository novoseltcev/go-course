package main

import (
	"github.com/caarlos0/env/v10"
	"github.com/spf13/cobra"

	"github.com/novoseltcev/go-course/internal/server"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Use this command to run server",
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			address, _ := flags.GetString("a")
			storeInterval, _ := flags.GetInt("s")
			fileStoragePath, _ := flags.GetString("f")
			restore, _ := flags.GetBool("r")
			databaseDsn, _ := flags.GetString("d")

			config := server.Config{
				Address: address,
				StoreInterval: storeInterval,
				FileStoragePath: fileStoragePath,
				Restore: restore,
				DatabaseDsn: databaseDsn,
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
	flags.StringP("d", "d", "", "Database connection string")
	flags.IntP("s", "s", 300, "backup interval")
	flags.StringP("f", "f", "/tmp/metrics-db.json", "path to backup")
	flags.BoolP("r", "r", true, "restore from backup after restart")
	return cmd
}

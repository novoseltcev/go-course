package main

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/novoseltcev/go-course/internal/server"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Use this command to run server",
		Run: func(cmd *cobra.Command, _ []string) {
			config, err := getConfig(cmd.Flags())
			if err != nil {
				log.Fatal(err)
			}

			s := server.NewServer(config)
			if err := s.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringP("a", "a", ":8080", "Server address")
	flags.StringP("d", "d", "", "Database connection string")
	flags.Int8P("s", "s", 0, "backup interval")
	flags.StringP("f", "f", "/tmp/metrics-db.json", "path to backup")
	flags.BoolP("r", "r", true, "restore from backup after restart")
	flags.StringP("k", "k", "", "Secret key for hashing data")

	return cmd
}

func getConfig(flags *pflag.FlagSet) (*server.Config, error) {
	address, err := flags.GetString("a")
	if err != nil {
		return nil, err
	}

	storeInterval, err := flags.GetInt8("s")
	if err != nil {
		return nil, err
	}

	fileStoragePath, err := flags.GetString("f")
	if err != nil {
		return nil, err
	}

	restore, err := flags.GetBool("r")
	if err != nil {
		return nil, err
	}

	databaseDsn, err := flags.GetString("d")
	if err != nil {
		return nil, err
	}

	secretKey, err := flags.GetString("k")
	if err != nil {
		return nil, err
	}

	config := server.Config{
		Address:         address,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
		DatabaseDsn:     databaseDsn,
		SecretKey:       secretKey,
	}

	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/novoseltcev/go-course/internal/server"
)

func Cmd() *cobra.Command {
	var configFile string

	cfg := &server.Config{} //nolint:exhaustruct
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Use this command to run server",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := parseConfig(cfg, configFile, cmd.Flags()); err != nil {
				log.Fatal(err)
			}

			s := server.NewServer(cfg)
			if err := s.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&configFile, "config", "c", "", "Path to config file")
	flags.StringVarP(&cfg.Address, "a", "a", cfg.Address, "Server address")
	flags.StringVarP(&cfg.DatabaseDsn, "d", "d", cfg.DatabaseDsn, "Database connection string")
	flags.StringVarP(&cfg.RawStoreInterval, "s", "s", cfg.RawStoreInterval, "Store interval")
	flags.StringVarP(&cfg.FileStoragePath, "f", "f", cfg.FileStoragePath, "Path to backup")
	flags.BoolVarP(&cfg.Restore, "r", "r", cfg.Restore, "Restore from backup after restart")
	flags.StringVarP(&cfg.SecretKey, "k", "k", cfg.SecretKey, "Secret key for hashing data")
	flags.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "Path to private key for decrypt data")

	return cmd
}

func parseConfig(cfg *server.Config, path string, flags *pflag.FlagSet) error {
	if path != "" {
		fd, err := os.Open(path)
		if err != nil {
			return err
		}

		defer fd.Close()

		if err := json.NewDecoder(fd).Decode(cfg); err != nil {
			return err
		}
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}

	if err := env.Parse(cfg); err != nil {
		return err
	}

	if err := cfg.FinishParse(); err != nil {
		return err
	}

	return nil
}

package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/novoseltcev/go-course/internal/agent"
)

func Cmd() *cobra.Command {
	var configFile string

	cfg := &agent.Config{}
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Use this command to run agent",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := parseConfig(cfg, configFile, cmd.Flags()); err != nil {
				log.Fatal(err)
			}

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			agent.Run(cfg, sigCh)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&configFile, "config", "c", "", "Path to config file")
	flags.StringVarP(&cfg.Address, "a", "a", cfg.Address, "Server address")
	flags.StringVarP(&cfg.RawPollInterval, "p", "p", cfg.RawPollInterval, "Poll runtime metrics interval")
	flags.StringVarP(&cfg.RawReportInterval, "r", "r", cfg.RawReportInterval, "Rate limit to send metrics")
	flags.StringVarP(&cfg.SecretKey, "k", "k", cfg.SecretKey, "Secret key for hashing data")
	flags.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "Path to public key to encrypt data")

	return cmd
}

func parseConfig(cfg *agent.Config, path string, flags *pflag.FlagSet) error {
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

package main

import (
	"compress/gzip"
	"crypto/sha256"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/pkg/chunkedrsa"
	"github.com/novoseltcev/go-course/pkg/compress"
	"github.com/novoseltcev/go-course/pkg/hash"
	"github.com/novoseltcev/go-course/pkg/retry"
)

// nolint: gochecknoglobals
var configFile string

// nolint: funlen
func Cmd() *cobra.Command {
	cfg := &agent.Config{
		Address:        "http://localhost:8080",
		PollInterval:   time.Second * 1,
		ReportInterval: time.Second * 5, // nolint:mnd
	}

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "CLI metrics agent",
		Run: func(cmd *cobra.Command, _ []string) {
			logger := logrus.New()
			fs := afero.NewOsFs()

			if err := cfg.Load(fs, configFile, cmd.Flags()); err != nil {
				logger.WithError(err).Panic("failed to parse config")
			}

			logger.SetLevel(logrus.DebugLevel) // TODO: set log level from config

			compressor, err := compress.NewGzip(gzip.BestCompression)
			if err != nil {
				logger.WithError(err).Panic("failed to initialize compressor")
			}

			opts := []reporters.Option{
				reporters.WithCompression(compressor),
				reporters.WithRetry(retry.Options{
					Retries:  3, // nolint:mnd
					Attempts: []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
				}),
			}
			if cfg.CryptoKey != "" {
				km, err := chunkedrsa.NewKeyManagerFromFile(fs, cfg.CryptoKey)
				if err != nil {
					logger.WithError(err).Panic("failed to load crypto key")
				}

				encryptor, err := chunkedrsa.NewEncryptor(km.MustPublicKey(), chunkedrsa.PKCS1v15Algorithm)
				if err != nil {
					logger.WithError(err).Panic("failed to initialize encryptor")
				}

				opts = append(opts, reporters.WithEncryption(encryptor))
			}

			if cfg.SecretKey != "" {
				opts = append(opts, reporters.WithCheckSum(hash.NewHMAC(cfg.SecretKey, sha256.New)))
			}

			reporter := reporters.NewHTTPClient(http.DefaultClient, cfg.Address, opts...)
			app := agent.NewApp(cfg, logger, reporter)

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer signal.Stop(sigCh)

			app.Run(sigCh)
		},
	}
	initFlags(cfg, cmd.Flags())

	return cmd
}

// initFlags initializes flags for parsing and help command.
func initFlags(cfg *agent.Config, flags *pflag.FlagSet) {
	flags.StringVarP(&configFile, "config", "c", "", "Path to config file")
	flags.StringVarP(&cfg.Address, "a", "a", cfg.Address, "Server address")
	flags.StringVarP(&cfg.RawPollInterval, "p", "p", cfg.RawPollInterval, "Poll runtime metrics interval")
	flags.StringVarP(&cfg.RawReportInterval, "r", "r", cfg.RawReportInterval, "Rate limit to send metrics")
	flags.StringVarP(&cfg.SecretKey, "k", "k", cfg.SecretKey, "Secret key for hashing data")
	flags.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "Path to public key to encrypt data")
}

package main

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/pkg/chunkedrsa"
	"github.com/novoseltcev/go-course/pkg/compress"
	"github.com/novoseltcev/go-course/pkg/hash"
	"github.com/novoseltcev/go-course/pkg/retry"
	pb "github.com/novoseltcev/go-course/proto/metrics"
)

// nolint: gochecknoglobals
var configFile string

// nolint: funlen
func Cmd() *cobra.Command {
	cfg := &agent.Config{
		Address:        "http://localhost:8080",
		PollInterval:   time.Second * 1,
		ReportInterval: time.Second * 5, // nolint:mnd
		ReporterType:   "http",
	}

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "CLI metrics agent",
		Run: func(cmd *cobra.Command, _ []string) {
			logger := logrus.New()
			fs := afero.NewOsFs()

			if err := cfg.Load(fs, configFile, cmd.Flags(), os.Args[1:]); err != nil {
				logger.WithError(err).Panic("failed to parse config")
			}

			logger.SetLevel(logrus.DebugLevel) // TODO: set log level from config

			var reporter agent.Reporter
			if cfg.ReporterType == "grpc" {
				client, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					logger.WithError(err).Panic("failed to dial gRPC server")
				}
				defer client.Close()

				reporter = reporters.NewGRPCReporter(pb.NewMetricsServiceClient(client))
			} else {
				var err error
				reporter, err = initHTTPReporter(cfg, fs)
				if err != nil {
					logger.WithError(err).Panic("failed to initialize reporter")
				}
			}

			app := agent.NewApp(cfg, logger, fs, reporter)

			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()

			app.Run(ctx)
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
	flags.StringVarP(&cfg.ReporterType, "reporter-type", "t", cfg.ReporterType, "Reporter type")
}

func initHTTPReporter(cfg *agent.Config, fs afero.Fs) (agent.Reporter, error) {
	compressor, err := compress.NewGzip(gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize compressor: %w", err)
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
			return nil, fmt.Errorf("failed to load crypto key: %w", err)
		}

		encryptor, err := chunkedrsa.NewEncryptor(km.MustPublicKey(), chunkedrsa.PKCS1v15Algorithm)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize encryptor: %w", err)
		}

		opts = append(opts, reporters.WithEncryption(encryptor))
	}

	if cfg.SecretKey != "" {
		opts = append(opts, reporters.WithCheckSum(hash.NewHMAC(cfg.SecretKey, sha256.New)))
	}

	return reporters.NewHTTPReporter(http.DefaultClient, cfg.Address, opts...), nil
}

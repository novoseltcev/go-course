package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/novoseltcev/go-course/internal/server"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/chunkedrsa"
)

// nolint: gochecknoglobals
var configFile string

// nolint: funlen
func Cmd() *cobra.Command {
	cfg := &server.Config{
		Address:         ":8080",
		Restore:         false,
		StoreInterval:   time.Second * 30, // nolint:mnd
		FileStoragePath: "/tmp/metrics-db.json",
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Metric server",
		Run: func(cmd *cobra.Command, _ []string) {
			logger := logrus.New()
			fs := afero.NewOsFs()

			if err := cfg.Load(fs, configFile, cmd.Flags(), os.Args[1:]); err != nil {
				logger.WithError(err).Fatal("failed to parse config")
			}

			logger.SetLevel(logrus.DebugLevel) // TODO: set log level from config

			var err error
			var db *sqlx.DB
			var storage storages.MetricStorager
			var decryptor chunkedrsa.Decryptor

			if cfg.DatabaseDsn == "" {
				storage = storages.NewMemStorage()
			} else {
				db, err = sqlx.Open("pgx", cfg.DatabaseDsn)
				if err != nil {
					logger.WithError(err).Panic("failed to open connection to database")
				}
				defer db.Close()

				storage = storages.NewPgStorage(db)
			}

			if cfg.CryptoKey != "" {
				km, err := chunkedrsa.NewKeyManagerFromFile(fs, cfg.CryptoKey)
				if err != nil {
					logger.WithError(err).Panic("failed to load crypto key")
				}

				decryptor, err = chunkedrsa.NewDecryptor(km.MustPrivateKey(), chunkedrsa.PKCS1v15Algorithm)
				if err != nil {
					logger.WithError(err).Panic("failed to initialize decryptor")
				}
			}

			app := server.NewApp(cfg, logger, fs, db, storage, decryptor)

			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()

			app.Run(ctx)
		},
	}
	initFlags(cfg, cmd.Flags())

	return cmd
}

// initFlags initializes flags for parsing and help command.
func initFlags(cfg *server.Config, flags *pflag.FlagSet) {
	flags.StringVarP(&configFile, "config", "c", "", "Path to config file")
	flags.StringVarP(&cfg.Address, "a", "a", cfg.Address, "Server address")
	flags.StringVarP(&cfg.DatabaseDsn, "d", "d", cfg.DatabaseDsn, "Database connection string")
	flags.StringVarP(&cfg.RawStoreInterval, "s", "s", cfg.RawStoreInterval, "Store interval")
	flags.StringVarP(&cfg.FileStoragePath, "f", "f", cfg.FileStoragePath, "Path to backup")
	flags.BoolVarP(&cfg.Restore, "r", "r", cfg.Restore, "Restore from backup after restart")
	flags.StringVarP(&cfg.SecretKey, "k", "k", cfg.SecretKey, "Secret key for hashing data")
	flags.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "Path to private key for decrypt data")
}

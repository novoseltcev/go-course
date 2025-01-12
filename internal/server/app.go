package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/chunkedrsa"
	"github.com/novoseltcev/go-course/pkg/httpserver"
	"github.com/novoseltcev/go-course/pkg/middlewares"
)

type App struct {
	cfg       *Config
	logger    *logrus.Logger
	db        *sqlx.DB
	storager  storages.MetricStorager
	decryptor chunkedrsa.Decryptor
}

func NewApp(
	cfg *Config,
	logger *logrus.Logger,
	db *sqlx.DB,
	storage storages.MetricStorager,
	decryptor chunkedrsa.Decryptor,
) *App {
	return &App{cfg: cfg, logger: logger, db: db, storager: storage, decryptor: decryptor}
}

func (app *App) Run(sigCh <-chan os.Signal) {
	ctx, cancel := context.WithCancel(context.Background())

	app.logger.Info("Server starting")

	if app.cfg.FileStoragePath != "" {
		defer Backup(app.cfg.FileStoragePath, app.storager) // nolint:errcheck

		if app.cfg.StoreInterval > 0 {
			go BackupWorker(ctx, app.cfg.StoreInterval, app.cfg.FileStoragePath, app.storager)
		}
	}

	if err := app.restore(); err != nil {
		app.logger.WithError(err).Error("failed to restore metrics")
	}

	srv := httpserver.New(configureRouter(app), httpserver.WithAddr(app.cfg.Address))

	app.logger.Info("Server started")

	go func() {
		defer cancel()

		select {
		case sig := <-sigCh:
			app.logger.WithField("signal", sig).Info("Signal received")
		case err := <-srv.Notify():
			if !errors.Is(err, http.ErrServerClosed) {
				app.logger.WithError(err).Error("Failed to listen and serve")
			}
		}

		app.logger.Info("Shutting down")

		if err := srv.Shutdown(); err != nil {
			app.logger.WithError(err).Error("Failed to shutdown")
		}
	}()

	<-ctx.Done()
	app.logger.Info("Server stopped")
}

func configureRouter(app *App) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Logger)

	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if app.db != nil {
			if err := app.db.PingContext(r.Context()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		}
	})

	r.Group(func(r chi.Router) {
		if app.cfg.SecretKey != "" {
			r.Use(middlewares.CheckSum(app.cfg.SecretKey))
		}

		r.Use(middlewares.Gzip)
		r.Use(middlewares.Decrypt(app.decryptor))

		if app.cfg.SecretKey != "" {
			r.Use(middlewares.Sign(app.cfg.SecretKey))
		}

		r.Mount("/", endpoints.NewAPIRouter(app.storager))
	})

	return r
}

func (app *App) restore() error {
	if !app.cfg.Restore || app.cfg.FileStoragePath == "" {
		return nil
	}

	fd, err := os.OpenFile(app.cfg.FileStoragePath, os.O_RDONLY, 0o666)
	if os.IsNotExist(err) {
		_, createErr := os.OpenFile(app.cfg.FileStoragePath, os.O_CREATE, 0o666)

		return errors.Join(err, createErr)
	}

	if err != nil {
		return err
	}

	defer fd.Close()

	return json.NewDecoder(fd).Decode(app.storager)
}

package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/httpserver"
	"github.com/novoseltcev/go-course/pkg/middlewares"
)

func Run(cfg *Config, sigCh <-chan os.Signal) { //nolint:cyclop
	logrus.SetLevel(logrus.InfoLevel) // TODO: set log level from config
	log := logrus.New().WithField("source", "server")
	doneCh := make(chan struct{}, 1)

	container, err := NewAppContainer(cfg)
	if err != nil {
		log.WithError(err).Fatal("failed to init app container")
	}
	defer container.Close()

	if cfg.FileStoragePath != "" {
		defer backup(cfg.FileStoragePath, container.Storage) //nolint:errcheck
	}

	log.Info("Server starting")

	if cfg.StoreInterval > 0 && cfg.FileStoragePath != "" {
		go func() {
			for {
				select {
				case <-doneCh:
					return
				default:
					time.Sleep(cfg.StoreInterval)
				}

				if err := backup(cfg.FileStoragePath, container.Storage); err != nil {
					log.WithError(err).Error("failed to backup metrics")
				}
			}
		}()
	}

	srv := httpserver.New(ConfigureRouter(container), httpserver.WithAddr(cfg.Address))

	log.Info("Server started")

	go func() {
		select {
		case sig := <-sigCh:
			log.WithField("signal", sig).Info("Signal received")
		case err := <-srv.Notify():
			if !errors.Is(err, http.ErrServerClosed) {
				log.WithError(err).Error("Failed to listen and serve")
			}
		}

		log.Info("Shutting down")

		if err := srv.Shutdown(); err != nil {
			log.WithError(err).Error("Failed to shutdown")
		}

		close(doneCh)
	}()

	<-doneCh
	log.Info("Server stopped")
}

func ConfigureRouter(container *AppContainer) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Logger)

	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	r.Group(func(r chi.Router) {
		if container.Cfg.SecretKey != "" {
			r.Use(middlewares.CheckSum(container.Cfg.SecretKey))
		}

		r.Use(middlewares.Gzip)
		r.Use(middlewares.Decrypt(container.Decryptor))

		if container.Cfg.SecretKey != "" {
			r.Use(middlewares.Sign(container.Cfg.SecretKey))
		}

		r.Mount("/", endpoints.NewAPIRouter(container.Storage))
	})

	return r
}

func restore(path string, storager storages.MetricStorager) error {
	fd, err := os.OpenFile(path, os.O_RDONLY, 0o666)
	if os.IsNotExist(err) {
		_, createErr := os.OpenFile(path, os.O_CREATE, 0o666)

		return errors.Join(err, createErr)
	}

	if err != nil {
		return err
	}

	defer fd.Close()

	return json.NewDecoder(fd).Decode(storager)
}

func backup(path string, storager storages.MetricStorager) error {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}

	defer fd.Close()

	if err := json.NewEncoder(fd).Encode(storager); err != nil {
		return err
	}

	return nil
}

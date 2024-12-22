package agent

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"

	"github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/agent/collectors"
	"github.com/novoseltcev/go-course/pkg/httpserver"
	"github.com/novoseltcev/go-course/pkg/workers"
)

func Run(cfg *Config, sigCh <-chan os.Signal) {
	logrus.SetLevel(logrus.InfoLevel) // TODO: set log level from config
	log := logrus.New().WithField("source", "agent")
	ctx, cancel := context.WithCancel(context.Background())

	container, err := NewAppContainer(cfg)
	if err != nil {
		log.WithError(err).Fatal("failed to init app container")
	}
	defer container.Close()

	log.Info("Agent starting")

	runtimeMetricCh := workers.Producer(ctx, collectors.CollectRuntimeMetrics, cfg.PollInterval)
	coreMetricCh := workers.Producer(ctx, collectors.CollectCoreMetrics, cfg.PollInterval)

	metricCh := workers.FanIn(ctx, runtimeMetricCh, coreMetricCh)

	go workers.AntiFraudConsumer(ctx, metricCh, container.Reporter.Report, cfg.ReportInterval)

	srv := httpserver.New(nil, httpserver.WithAddr(":9000"))

	log.Info("Agent started")

	go func() {
		defer cancel()

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
	}()

	<-ctx.Done()
	log.Info("Agent stopped")
}

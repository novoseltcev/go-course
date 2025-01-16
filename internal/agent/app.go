package agent

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof" // nolint:gosec
	"os"

	"github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/agent/collectors"
	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/pkg/httpserver"
	"github.com/novoseltcev/go-course/pkg/workers"
)

type App struct {
	cfg      *Config
	logger   *logrus.Logger
	reporter reporters.Reporter
}

func NewApp(cfg *Config, logger *logrus.Logger, reporter reporters.Reporter) *App {
	return &App{cfg: cfg, logger: logger, reporter: reporter}
}

func (app *App) Run(sigCh <-chan os.Signal) {
	ctx, cancel := context.WithCancel(context.Background())

	app.logger.Info("Agent starting")

	runtimeMetricCh := workers.Producer(ctx, collectors.CollectRuntimeMetrics, app.cfg.PollInterval)
	coreMetricCh := workers.Producer(ctx, collectors.CollectCoreMetrics, app.cfg.PollInterval)

	metricCh := workers.FanIn(ctx, runtimeMetricCh, coreMetricCh)

	go workers.AntiFraudConsumer(ctx, metricCh, app.reporter.Report, app.cfg.ReportInterval)

	srv := httpserver.New(nil, httpserver.WithAddr(":9000"))

	app.logger.Info("Agent started")

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
	app.logger.Info("Agent stopped")
}

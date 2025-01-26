package agent

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof" // nolint:gosec
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/novoseltcev/go-course/internal/agent/collectors"
	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/pkg/httpserver"
	"github.com/novoseltcev/go-course/pkg/workers"
)

const shutdownTimeout = 5 * time.Second

//go:generate mockgen -source=app.go -destination=./app_mock_test.go -package=agent_test -typed
type Reporter interface {
	Report(ctx context.Context, metrics []schemas.Metric) error
}

type App struct {
	cfg      *Config
	logger   *logrus.Logger
	fs       afero.Fs
	reporter Reporter
}

func NewApp(cfg *Config, logger *logrus.Logger, fs afero.Fs, reporter Reporter) *App {
	return &App{cfg: cfg, logger: logger, fs: fs, reporter: reporter}
}

func (app *App) Run(ctx context.Context) {
	app.logger.Info("Agent starting")

	runtimeMetricCh := workers.Producer(ctx, collectors.CollectRuntimeMetrics, app.cfg.PollInterval)
	coreMetricCh := workers.Producer(ctx, collectors.CollectCoreMetrics, app.cfg.PollInterval)

	metricCh := workers.FanIn(ctx, runtimeMetricCh, coreMetricCh)

	go workers.AntiFraudConsumer(ctx, metricCh, app.reporter.Report, app.cfg.ReportInterval)

	srv := httpserver.New(nil, httpserver.WithAddr(":9000"))
	go srv.Run()

	app.logger.Info("Agent started")

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)

		select {
		case <-ctx.Done():
			app.logger.Info("Run is interrupted by context")
		case err := <-srv.Notify():
			if !errors.Is(err, http.ErrServerClosed) {
				app.logger.WithError(err).Error("Failed to listen and serve")
			}
		}

		app.logger.Info("Shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil { // nolint: contextcheck
			app.logger.WithError(err).Error("Failed to shutdown")
		}
	}()

	<-doneCh
	app.logger.Info("Agent stopped")
}

package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/server"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

func TestRun(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	app := server.NewApp(&server.Config{}, logrus.New(), afero.NewMemMapFs(), nil, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	app.Run(ctx)
}

const testFile = "test.json"

func TestRunWithRestoreFailsFileNotExists(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	app := server.NewApp(&server.Config{
		SecretKey:       "secret",
		Restore:         true,
		FileStoragePath: testFile,
		StoreInterval:   time.Second,
	}, logrus.New(), afero.NewMemMapFs(), nil, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	app.Run(ctx)
}

func TestRunWithRestoreFailsUnmarsalError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, testutils.Bytes)

	app := server.NewApp(&server.Config{
		SecretKey:       "secret",
		Restore:         true,
		FileStoragePath: testFile,
		StoreInterval:   time.Second,
	}, logrus.New(), fs, nil, storages.NewMemStorage(), nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	app.Run(ctx)
}

func TestRunWithRestoreSuccess(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, []byte(`{"gauge":[{"id": "test", "type": "gauge"}], "counter":[]}`))

	app := server.NewApp(&server.Config{
		SecretKey:       "secret",
		Restore:         true,
		FileStoragePath: testFile,
		StoreInterval:   time.Second,
	}, logrus.New(), fs, nil, storages.NewMemStorage(), nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	app.Run(ctx)
}

package server_test

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

func TestBackup_WithoutFile_Success(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	store := storages.NewMemStorage()

	testDelta := int64(10)
	testValue := 10.123
	err := store.SaveBatch(context.Background(), []schemas.Metric{
		{ID: "test", MType: "gauge", Value: &testValue},
		{ID: "test", MType: "counter", Delta: &testDelta},
	})
	require.NoError(t, err)

	require.NoError(t, server.Backup(fs, testFile, store))

	data, err := afero.ReadFile(fs, testFile)
	require.NoError(t, err)

	assert.JSONEq(t, `{"Data": {"gauge":{"test": {"value": 10.123}}, "counter":{"test": {"delta": 10}}}}`, string(data))
}

func TestBackup_WithFile_Success(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, testutils.Bytes)
	store := storages.NewMemStorage()

	testDelta := int64(10)
	testValue := 10.123
	err := store.SaveBatch(context.Background(), []schemas.Metric{
		{ID: "test", MType: "gauge", Value: &testValue},
		{ID: "test", MType: "counter", Delta: &testDelta},
	})
	require.NoError(t, err)

	require.NoError(t, server.Backup(fs, testFile, store))

	data, err := afero.ReadFile(fs, testFile)
	require.NoError(t, err)

	assert.JSONEq(t, `{"Data": {"gauge":{"test": {"value": 10.123}}, "counter":{"test": {"delta": 10}}}}`, string(data))
}

func TestBackupFailsOpen(t *testing.T) {
	t.Parallel()

	err := server.Backup(afero.NewReadOnlyFs(afero.NewMemMapFs()), testFile, nil)
	require.ErrorIs(t, err, os.ErrPermission)
}

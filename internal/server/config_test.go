package server_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/server"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

var flags = pflag.NewFlagSet("test", pflag.ContinueOnError)

func TestConfig_ParseRawFields(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		cfg := server.Config{}
		require.NoError(t, cfg.ParseRawFields())

		assert.Equal(t, time.Duration(0), cfg.StoreInterval)
	})

	t.Run("parse", func(t *testing.T) {
		t.Parallel()

		cfg := server.Config{RawStoreInterval: "1s"}
		require.NoError(t, cfg.ParseRawFields())

		assert.Equal(t, time.Second, cfg.StoreInterval)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		cfg := server.Config{RawStoreInterval: "1"}

		assert.ErrorContains(t, cfg.ParseRawFields(), "time: ")
	})
}

func TestConfig_Load_Success(t *testing.T) {
	t.Parallel()

	var cfg server.Config

	require.NoError(t,
		cfg.Load(afero.NewMemMapFs(), "", flags, nil),
	)

	assert.Equal(t, server.Config{}, cfg)
}

func TestConfig_Load_WithFile_Success(t *testing.T) {
	t.Parallel()

	var cfg server.Config

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, []byte(`
		{
			"address": "test",
			"restore": true,
			"store_file": "test.json",
			"store_interval": "1s",
			"database_dsn": "some",
			"crypto_key": "key"
		}`),
	)

	require.NoError(t, cfg.Load(fs, testFile, flags, nil))

	assert.Equal(t, server.Config{
		Address:          "test",
		Restore:          true,
		FileStoragePath:  "test.json",
		RawStoreInterval: "1s",
		StoreInterval:    time.Second,
		DatabaseDsn:      "some",
		CryptoKey:        "key",
	}, cfg)
}

func TestConfig_Load_FailedOpen(t *testing.T) {
	t.Parallel()

	var cfg server.Config

	require.ErrorIs(t, cfg.Load(afero.NewMemMapFs(), testFile, flags, nil), os.ErrNotExist)
}

func TestConfig_Load_FailedDecodeJSON(t *testing.T) {
	t.Parallel()

	var cfg server.Config

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, []byte(`{`))

	assert.ErrorIs(t, cfg.Load(fs, testFile, flags, nil), io.ErrUnexpectedEOF)
}

func TestConfig_Load_FailedParseRawFields(t *testing.T) {
	t.Parallel()

	var cfg server.Config

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, []byte(`{"store_interval": "1"}`))

	assert.ErrorContains(t, cfg.Load(fs, testFile, flags, nil), "time: ")
}

func TestConfig_Load_FailedParseFlags(t *testing.T) {
	t.Parallel()

	var cfg server.Config

	assert.ErrorContains(t,
		cfg.Load(afero.NewMemMapFs(), "", flags, []string{"---a", "test"}),
		"bad flag syntax: ",
	)
}

package agent_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/agent"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

var flags = pflag.NewFlagSet("test", pflag.ContinueOnError)

func TestConfig_ParseRawFields(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		cfg := agent.Config{}
		require.NoError(t, cfg.ParseRawFields())

		assert.Equal(t, time.Duration(0), cfg.PollInterval)
		assert.Equal(t, time.Duration(0), cfg.ReportInterval)
	})

	t.Run("parse", func(t *testing.T) {
		t.Parallel()

		cfg := agent.Config{
			RawPollInterval:   "1µs",
			RawReportInterval: "2s",
		}
		require.NoError(t, cfg.ParseRawFields())

		assert.Equal(t, time.Microsecond, cfg.PollInterval)
		assert.Equal(t, 2*time.Second, cfg.ReportInterval)
	})

	t.Run("error RawPollInterval", func(t *testing.T) {
		t.Parallel()

		cfg := agent.Config{
			RawPollInterval:   "1",
			RawReportInterval: "a",
		}

		assert.ErrorContains(t, cfg.ParseRawFields(), "time: ")
	})

	t.Run("error RawReportInterval", func(t *testing.T) {
		t.Parallel()

		cfg := agent.Config{
			RawPollInterval:   "1s",
			RawReportInterval: "a",
		}

		assert.ErrorContains(t, cfg.ParseRawFields(), "time: ")
	})
}

const testFile = "test.json"

func TestConfig_Load_WithoutFile_Success(t *testing.T) {
	t.Parallel()

	var cfg agent.Config

	require.NoError(t, cfg.Load(afero.NewMemMapFs(), "", flags, nil))

	assert.Equal(t, agent.Config{}, cfg)
}

func TestConfig_Load_WithFile_Success(t *testing.T) {
	t.Parallel()

	var cfg agent.Config

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, []byte(`
		{
			"address": "test",
			"poll_interval": "1s",
			"report_interval": "2s",
			"crypto_key": "key"
		}`),
	)

	require.NoError(t, cfg.Load(fs, testFile, flags, nil))

	assert.Equal(t, agent.Config{
		Address:           "test",
		CryptoKey:         "key",
		RawPollInterval:   "1s",
		RawReportInterval: "2s",
		PollInterval:      time.Second,
		ReportInterval:    2 * time.Second,
	}, cfg)
}

func TestConfig_Load_FailedOpen(t *testing.T) {
	t.Parallel()

	var cfg agent.Config

	require.ErrorIs(t, cfg.Load(afero.NewMemMapFs(), testFile, flags, nil), os.ErrNotExist)
}

func TestConfig_Load_FailedDecodeJSON(t *testing.T) {
	t.Parallel()

	var cfg agent.Config

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, []byte(`{`))

	assert.ErrorIs(t, cfg.Load(fs, testFile, flags, nil), io.ErrUnexpectedEOF)
}

func TestConfig_Load_FailedParseRawFields(t *testing.T) {
	t.Parallel()

	var cfg agent.Config

	fs := afero.NewMemMapFs()
	helpers.WriteToFile(t, fs, testFile, []byte(`{"poll_interval": "1"}`))

	assert.ErrorContains(t, cfg.Load(fs, testFile, flags, nil), "time: ")
}

func TestConfig_Load_FailedParseFlags(t *testing.T) {
	t.Parallel()

	var cfg agent.Config

	assert.ErrorContains(t,
		cfg.Load(afero.NewMemMapFs(), "", flags, []string{"---a", "test"}),
		"bad flag syntax: ",
	)
}

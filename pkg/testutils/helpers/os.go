package helpers

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func WriteToFile(t *testing.T, fs afero.Fs, path string, data []byte) {
	t.Helper()

	fd, err := fs.Create(path)
	require.NoError(t, err)

	_, err = fd.Write(data)
	require.NoError(t, err)

	err = fd.Close()
	require.NoError(t, err)
}

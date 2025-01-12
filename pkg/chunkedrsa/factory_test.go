package chunkedrsa_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa"
	"github.com/novoseltcev/go-course/pkg/chunkedrsa/pkcs1v15"
)

func TestConsts(t *testing.T) {
	t.Parallel()

	require.Equal(t, "PKCS1v15", chunkedrsa.PKCS1v15Algorithm)
}

func TestNewEncryptor(t *testing.T) {
	t.Parallel()

	enc, err := chunkedrsa.NewEncryptor(nil, "PKCS1v15")
	require.NoError(t, err)
	assert.IsType(t, &pkcs1v15.Encryptor{}, enc)
}

func TestNewEncryptorFailWithUnsupportedAlgorithm(t *testing.T) {
	t.Parallel()

	for _, algorithm := range []string{"pkcs1v15", "OAEP", "UNKNOWN"} {
		t.Run(algorithm, func(t *testing.T) {
			t.Parallel()

			enc, err := chunkedrsa.NewEncryptor(nil, algorithm)
			require.Nil(t, enc)
			assert.ErrorIs(t, err, chunkedrsa.ErrUnsupportedAlgorithm)
		})
	}
}

func TestNewDecryptor(t *testing.T) {
	t.Parallel()

	dec, err := chunkedrsa.NewDecryptor(nil, chunkedrsa.PKCS1v15Algorithm)
	require.NoError(t, err)
	assert.IsType(t, &pkcs1v15.Decryptor{}, dec)
}

func TestNewDecryptorFailWithUnsupportedAlgorithm(t *testing.T) {
	t.Parallel()

	for _, algorithm := range []string{"pkcs1v15", "OAEP", "UNKNOWN"} {
		t.Run(algorithm, func(t *testing.T) {
			t.Parallel()

			dec, err := chunkedrsa.NewDecryptor(nil, algorithm)
			require.Nil(t, dec)
			assert.ErrorIs(t, err, chunkedrsa.ErrUnsupportedAlgorithm)
		})
	}
}

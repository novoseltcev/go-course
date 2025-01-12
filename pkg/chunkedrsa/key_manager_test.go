package chunkedrsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa"
)

func TestMustPanics(t *testing.T) {
	t.Parallel()

	km := chunkedrsa.NewPublicKeyManager(nil)

	t.Run("public key", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithValue(t, "public key is not set", func() { _ = km.MustPublicKey() })
	})

	t.Run("private key", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithValue(t, "private key is not set", func() { _ = km.MustPrivateKey() })
	})
}

func TestNewKeyManager(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	km := chunkedrsa.NewKeyManager(key)
	assert.Equal(t, key, km.MustPrivateKey())
	assert.Equal(t, &key.PublicKey, km.MustPublicKey())
}

func TestNewPublicKeyManager(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	km := chunkedrsa.NewPublicKeyManager(&key.PublicKey)
	assert.Equal(t, &key.PublicKey, km.MustPublicKey())
}

const testFile = "test.pem"

func TestNewKeyManagerFromFileSuccessPrivateKey(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(
		fs,
		testFile,
		pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}),
		0o644,
	))

	km, err := chunkedrsa.NewKeyManagerFromFile(fs, testFile)
	require.NoError(t, err)

	assert.Equal(t, key, km.MustPrivateKey())
	assert.Equal(t, &key.PublicKey, km.MustPublicKey())
}

func TestNewKeyManagerFromFileSuccessPublicKey(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	fs := afero.NewMemMapFs()

	require.NoError(t, afero.WriteFile(
		fs,
		testFile,
		pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey)}),
		0o644,
	))

	km, err := chunkedrsa.NewKeyManagerFromFile(fs, testFile)
	require.NoError(t, err)

	assert.Equal(t, &key.PublicKey, km.MustPublicKey())
	assert.Panics(t, func() { _ = km.MustPrivateKey() })
}

func TestNewKeyManagerFromFileFailWithNotExistFile(t *testing.T) {
	t.Parallel()

	_, err := chunkedrsa.NewKeyManagerFromFile(afero.NewMemMapFs(), testFile)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestNewKeyManagerFromFileFailWithBadPemData(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	keyPemData := []byte("bad pem data")

	require.NoError(t, afero.WriteFile(fs, testFile, keyPemData, 0o644))

	_, err := chunkedrsa.NewKeyManagerFromFile(fs, testFile)
	assert.ErrorIs(t, err, chunkedrsa.ErrBadPemData)
}

func TestNewKeyManagerFromFileFailWithUnknownPemType(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()

	require.NoError(t, afero.WriteFile(
		fs, testFile, pem.EncodeToMemory(&pem.Block{Type: "UNKNOWN", Bytes: nil}), 0o644,
	))

	_, err := chunkedrsa.NewKeyManagerFromFile(fs, testFile)
	assert.ErrorIs(t, err, chunkedrsa.ErrBadPemType)
}

func TestNewKeyManagerFromFileFailParsePrivateKey(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(
		fs, testFile, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: nil}), 0o644,
	))

	_, err := chunkedrsa.NewKeyManagerFromFile(fs, testFile)
	assert.Error(t, err)
}

func TestNewKeyManagerFromFileFailParsePublicKey(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()

	require.NoError(t, afero.WriteFile(
		fs, testFile, pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: nil}), 0o644,
	))

	_, err := chunkedrsa.NewKeyManagerFromFile(fs, testFile)
	assert.Error(t, err)
}

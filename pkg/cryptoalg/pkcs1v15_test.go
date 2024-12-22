package cryptoalg_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/cryptoalg"
)

func TestEncrypt(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	encryptor := cryptoalg.NewPKCS1v15Manager(key)

	data := []byte("test")
	encrypted, err := encryptor.Encrypt(data)
	require.NoError(t, err)

	decrypted, err := key.Decrypt(rand.Reader, encrypted, nil)
	require.NoError(t, err)
	assert.Equal(t, data, decrypted)
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	encryptor := cryptoalg.NewPKCS1v15Manager(key)

	data := []byte("test")
	encrypted, err := encryptor.Encrypt(data)
	require.NoError(t, err)

	decrypted, err := encryptor.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, data, decrypted)
}

func TestEncryptChunks(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 128)
	require.NoError(t, err)

	mng := cryptoalg.NewPKCS1v15Manager(key)

	data := []byte("testtt")
	require.Greater(t, len(data), key.Size()-11)

	encrypted, err := mng.Encrypt(data)
	require.NoError(t, err)

	decrypted, err := mng.Decrypt(encrypted)
	require.NoError(t, err)
	assert.Equal(t, data, decrypted)
}

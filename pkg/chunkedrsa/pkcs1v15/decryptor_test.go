package pkcs1v15_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa/pkcs1v15"
)

func encrypt(t *testing.T, key *rsa.PublicKey, data []byte) []byte {
	t.Helper()
	require.LessOrEqual(t, len(data), key.Size()-11)

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, key, data)
	require.NoError(t, err)

	return encrypted
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	dec := pkcs1v15.NewDecryptor(key)

	data := []byte("test")
	encrypted := encrypt(t, &key.PublicKey, data)

	decrypted, err := dec.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, data, decrypted)
}

func TestDecryptChunks(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 128)
	require.NoError(t, err)
	dec := pkcs1v15.NewDecryptor(key)

	data := []byte("testtt")
	chunkSize := key.Size() - 11
	require.Greater(t, len(data), chunkSize)
	require.LessOrEqual(t, len(data), 2*chunkSize)

	encrypted := encrypt(t, &key.PublicKey, data[:key.Size()-11])
	encrypted = append(encrypted, encrypt(t, &key.PublicKey, data[key.Size()-11:])...)

	decrypted, err := dec.Decrypt(encrypted)
	require.NoError(t, err)

	assert.Equal(t, data, decrypted)
}

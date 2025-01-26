package pkcs1v15_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa/pkcs1v15"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

func decrypt(t *testing.T, key *rsa.PrivateKey, data []byte) []byte {
	t.Helper()
	require.LessOrEqual(t, len(data), key.Size())

	decrypted, err := key.Decrypt(rand.Reader, data, nil)
	require.NoError(t, err)

	return decrypted
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	enc := pkcs1v15.NewEncryptor(&key.PublicKey)

	data := []byte("test")
	encrypted, err := enc.Encrypt(data)
	require.NoError(t, err)

	assert.Equal(t, data, decrypt(t, key, encrypted))
}

func TestEncryptChunks(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 128)
	require.NoError(t, err)
	enc := pkcs1v15.NewEncryptor(&key.PublicKey)

	data := []byte("testtt")
	encrypted, err := enc.Encrypt(data)
	require.NoError(t, err)

	require.Greater(t, len(encrypted), key.Size())
	require.LessOrEqual(t, len(encrypted), 2*key.Size())

	decrypted := decrypt(t, key, encrypted[:key.Size()])
	decrypted = append(decrypted, decrypt(t, key, encrypted[key.Size():])...)

	assert.Equal(t, data, decrypted)
}

func TestEncrypt_FailsChunking(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 88) // >=89 has non-zero data padding
	require.NoError(t, err)

	enc := pkcs1v15.NewEncryptor(&key.PublicKey)

	_, err = enc.Encrypt(testutils.Bytes)
	assert.Error(t, err)
}

func TestEncrypt_FailsChunkEncryption(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	key.PublicKey.E = 0 // Make public key to invalid

	enc := pkcs1v15.NewEncryptor(&key.PublicKey)

	_, err = enc.Encrypt(testutils.Bytes)
	assert.Error(t, err, "crypto/rsa: public exponent too small")
}

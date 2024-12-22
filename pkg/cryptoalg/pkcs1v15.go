package cryptoalg

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"
)

var ErrBadPemData = errors.New("bad key data: not PEM-encoded")

type PKCS1v15Algorithm struct {
	key       *rsa.PrivateKey
	publicKey *rsa.PublicKey
}

func NewPKCS1v15Manager(key *rsa.PrivateKey) *PKCS1v15Algorithm {
	return &PKCS1v15Algorithm{key: key, publicKey: &key.PublicKey}
}

func NewPKCS1v15DecryptorFromFile(path string) (*PKCS1v15Algorithm, error) {
	pemBlock, err := loadPemBlock(path)
	if err != nil {
		return nil, err
	}

	if pemBlock.Type != "RSA PRIVATE KEY" {
		return nil, ErrBadPemData
	}

	key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &PKCS1v15Algorithm{key: key, publicKey: nil}, nil
}

func NewPKCS1v15EncryptorFromFile(path string) (*PKCS1v15Algorithm, error) {
	pemBlock, err := loadPemBlock(path)
	if err != nil {
		return nil, err
	}

	if pemBlock.Type != "RSA PUBLIC KEY" {
		return nil, ErrBadPemData
	}

	key, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &PKCS1v15Algorithm{publicKey: key, key: nil}, nil
}

func loadPemBlock(path string) (*pem.Block, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	b, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, ErrBadPemData
	}

	return block, nil
}

// Encrypt encrypts data with PKCS1v15 padding.
//
// Data is split into chunks of the size of the key and 11 bytes for PKCS1 v1.5 padding.
func (alg *PKCS1v15Algorithm) Encrypt(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	chunks := splitToChunks(data, alg.publicKey.Size()-11) //nolint: mnd // 11 bytes for PKCS1 v1.5 padding
	for _, chunk := range chunks {
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, alg.publicKey, chunk)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(encryptedChunk); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Decrypt decrypts data with PKCS1v15 padding.
//
// Data is split into chunks of the size of the key. 11 bytes for PKCS1 v1.5 padding is included in the chunk size.
func (alg *PKCS1v15Algorithm) Decrypt(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	chunks := splitToChunks(data, alg.key.Size()) // 11 bytes for PKCS1 v1.5 padding
	for _, chunk := range chunks {
		decryptedChunk, err := rsa.DecryptPKCS1v15(rand.Reader, alg.key, chunk)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(decryptedChunk); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// splitToChunks splits data into chunks of the size of chunkSize.
//
// >>> splitToChunks([]int{1, 2, 3, 4, 5}, 2)
// [][]int{{1, 2}, {3, 4}, {5}}.
func splitToChunks[T any](data []T, chunkSize int) [][]T {
	chunks := make([][]T, 0)

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunks = append(chunks, data[i:end])
	}

	return chunks
}

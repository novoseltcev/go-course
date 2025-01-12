// Package chunkedrsa provides chunked encryption and decryption of data using RSA algorithms with.
package chunkedrsa

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa/pkcs1v15"
)

// Encryptor is an interface for encrypting data.
type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
}

// Decryptor is an interface for decrypting data.
type Decryptor interface {
	Decrypt(data []byte) ([]byte, error)
}

const (
	PKCS1v15Algorithm = "PKCS1v15"
)

var ErrUnsupportedAlgorithm = errors.New("unsupported algorithm")

// NewEncryptor creates an Encryptor instance for the given key and algorithm.
//
// The algorithm must be one of the following:
//
//   - "PKCS1v15"
func NewEncryptor(key *rsa.PublicKey, algorithm string) (Encryptor, error) {
	switch algorithm {
	case PKCS1v15Algorithm:
		return pkcs1v15.NewEncryptor(key), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, algorithm)
	}
}

// NewDecryptor creates a Decryptor instance for the given key and algorithm.
//
// The algorithm must be one of the following:
//
//   - "PKCS1v15"
func NewDecryptor(key *rsa.PrivateKey, algorithm string) (Decryptor, error) {
	switch algorithm {
	case PKCS1v15Algorithm:
		return pkcs1v15.NewDecryptor(key), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, algorithm)
	}
}

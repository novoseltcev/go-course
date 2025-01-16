package chunkedrsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"

	"github.com/spf13/afero"
)

var (
	ErrBadPemData = errors.New("bad key data: not PEM-encoded")
	ErrBadPemType = errors.New("bad key data: unknown PEM type")
)

// KeyManager is a wrapper for RSA keys.
type KeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewKeyManager creates a new KeyManager instance for the given private key.
func NewKeyManager(key *rsa.PrivateKey) *KeyManager {
	return &KeyManager{privateKey: key, publicKey: &key.PublicKey}
}

// NewPublicKeyManager creates a new KeyManager instance for the given public key.
func NewPublicKeyManager(key *rsa.PublicKey) *KeyManager {
	return &KeyManager{publicKey: key}
}

// NewKeyManagerFromFile creates a new KeyManager instance by loading the key from the given pem-encoded file.
func NewKeyManagerFromFile(fs afero.Fs, path string) (*KeyManager, error) {
	pemBlock, err := loadPemBlock(fs, path)
	if err != nil {
		return nil, err
	}

	switch pemBlock.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}

		return &KeyManager{privateKey: key, publicKey: &key.PublicKey}, nil
	case "RSA PUBLIC KEY":
		key, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}

		return &KeyManager{privateKey: nil, publicKey: key}, nil
	default:
		return nil, ErrBadPemType
	}
}

// MustPrivateKey returns the private key or panics if it is not set.
func (m *KeyManager) MustPrivateKey() *rsa.PrivateKey {
	if m.privateKey == nil {
		panic("private key is not set")
	}

	return m.privateKey
}

// MustPublicKey returns the public key or panics if it is not set.
func (m *KeyManager) MustPublicKey() *rsa.PublicKey {
	if m.publicKey == nil {
		panic("public key is not set")
	}

	return m.publicKey
}

func loadPemBlock(fs afero.Fs, path string) (*pem.Block, error) {
	fd, err := fs.Open(path)
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

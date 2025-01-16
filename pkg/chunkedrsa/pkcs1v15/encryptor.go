package pkcs1v15

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa/utils"
)

type Encryptor struct {
	publicKey *rsa.PublicKey
}

func NewEncryptor(key *rsa.PublicKey) *Encryptor {
	return &Encryptor{publicKey: key}
}

// Encrypt encrypts data with PKCS1v15 padding.
//
// Data is split into chunks of the size of the key and 11 bytes for PKCS1 v1.5 padding.
func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	chunks := utils.SplitToChunks(data, e.publicKey.Size()-11) // nolint: mnd // 11 bytes for PKCS1 v1.5 padding
	for _, chunk := range chunks {
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, e.publicKey, chunk)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(encryptedChunk); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

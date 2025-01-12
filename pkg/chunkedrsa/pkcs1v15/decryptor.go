package pkcs1v15

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"

	"github.com/novoseltcev/go-course/pkg/chunkedrsa/utils"
)

type Decryptor struct {
	privateKey *rsa.PrivateKey
}

func NewDecryptor(key *rsa.PrivateKey) *Decryptor {
	return &Decryptor{privateKey: key}
}

// Decrypt decrypts data with PKCS1v15 padding.
//
// Data is split into chunks of the size of the key. 11 bytes for PKCS1 v1.5 padding is included in the chunk size.
func (d *Decryptor) Decrypt(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	chunks := utils.SplitToChunks(data, d.privateKey.Size()) // 11 bytes for PKCS1 v1.5 padding
	for _, chunk := range chunks {
		decryptedChunk, err := rsa.DecryptPKCS1v15(rand.Reader, d.privateKey, chunk)
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(decryptedChunk); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

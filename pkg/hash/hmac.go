package hash

import (
	"crypto/hmac"
	"hash"
)

type HMAC struct {
	fn  func() hash.Hash
	key []byte
}

func NewHMAC(key string, fn func() hash.Hash) *HMAC {
	return &HMAC{
		key: []byte(key),
		fn:  fn,
	}
}

func (hm *HMAC) GetHash(data []byte) ([]byte, error) {
	fn := hmac.New(hm.fn, hm.key)

	_, err := fn.Write(data)
	if err != nil {
		return nil, err
	}

	return fn.Sum(nil), nil
}

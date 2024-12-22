package hash_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/hash"
)

func TestGetHash(t *testing.T) {
	t.Parallel()

	hmac := hash.NewHMAC("secret", sha256.New)
	hash, err := hmac.GetHash([]byte("test"))
	require.NoError(t, err)

	assert.Equal(t,
		// computed by site: https://www.devglan.com/online-tools/hmac-sha256-online
		"0329a06b62cd16b33eb6792be8c60b158d89a2ee3a876fce9a881ebb488c0914",
		hex.EncodeToString(hash),
	)
}

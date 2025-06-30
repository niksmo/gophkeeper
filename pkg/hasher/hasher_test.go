package hasher_test

import (
	"crypto/rand"
	"testing"

	"github.com/niksmo/gophkeeper/pkg/hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCryptoHasher(t *testing.T) {
	const cost = 10

	h := hasher.NewCryptoHasher(cost)

	t.Run("Generate", func(t *testing.T) {
		t.Run("Ordinary", func(t *testing.T) {
			password := []byte("simplePassword")
			hash, err := h.Generate(password)
			require.NoError(t, err)
			assert.NotZero(t, hash)
		})
		t.Run("TooLongSource", func(t *testing.T) {
			const passwordLen = 80
			password := make([]byte, passwordLen)
			n, err := rand.Read(password)
			require.NoError(t, err)
			require.Equal(t, passwordLen, n)
			hash, err := h.Generate(password)
			require.ErrorIs(t, err, hasher.ErrTooLongSrc)
			assert.Zero(t, hash)
		})
	})

	t.Run("Compare", func(t *testing.T) {
		t.Run("Ordinary", func(t *testing.T) {
			password := []byte("simplePassword")
			hash, err := h.Generate(password)
			require.NoError(t, err)
			assert.NotZero(t, hash)

			err = h.Compare(hash, password)
			require.NoError(t, err)
		})

		t.Run("NotMatch", func(t *testing.T) {
			password := []byte("simplePassword")
			hash, err := h.Generate(password)
			require.NoError(t, err)
			assert.NotZero(t, hash)

			password = []byte("changedPassword")
			err = h.Compare(hash, password)
			require.Error(t, err)
		})

		t.Run("NotHash", func(t *testing.T) {
			password := []byte("simplePassword")
			hash := []byte("notHash")
			err := h.Compare(hash, password)
			require.Error(t, err)
		})
	})

}

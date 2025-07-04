package cipher_test

import (
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/niksmo/gophkeeper/pkg/cipher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	alphabet = []rune(`1234567890-=!@#$%^&*()_+[]\;',./{}|:"<>?` +
		`qwertyuiopasdfghjklzxcvbnm` +
		`QWERTYUIOPASDFGHJKLZXCVBNM`)
	alen = len(alphabet)
)

func getRandPwd(len int) string {
	var b strings.Builder
	for range len {
		b.WriteRune(alphabet[rand.IntN(alen)])
	}
	return b.String()
}

func TestEncrypter(t *testing.T) {
	data := []byte("hello_world")

	e := cipher.NewEncrypter()
	e.SetKey(getRandPwd(100))
	encryptedData, _ := e.Encrypt(data)
	require.NotEqual(t, data, encryptedData)
}

func TestDecrypter(t *testing.T) {
	t.Run("ValidPassword", func(t *testing.T) {
		password := getRandPwd(100)
		data := []byte("hello_world")

		e := cipher.NewEncrypter()
		e.SetKey(password)
		encryptedData, _ := e.Encrypt(data)
		require.NotEqual(t, data, encryptedData)

		d := cipher.NewDecrypter()
		d.SetKey(password)
		decryptedData, err := d.Decrypt(encryptedData)
		require.NoError(t, err)
		assert.Equal(t, data, decryptedData)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		data := []byte("hello_world")

		e := cipher.NewEncrypter()
		e.SetKey(getRandPwd(100))
		encryptedData, _ := e.Encrypt(data)
		require.NotEqual(t, data, encryptedData)

		d := cipher.NewDecrypter()
		d.SetKey(getRandPwd(100))
		decryptedData, err := d.Decrypt(encryptedData)
		require.Error(t, err)
		assert.NotEqual(t, data, decryptedData)
	})
}

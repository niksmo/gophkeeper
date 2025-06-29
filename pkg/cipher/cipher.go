package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

const (
	keySize = 32
	iter    = 4096
)

type Encrypter struct {
	keySetter
}

func NewEncrypter() *Encrypter {
	return &Encrypter{}
}

func (e *Encrypter) Encrypt(data []byte) ([]byte, error) {
	const op = "Encrypter.Encrypt"

	key, err := makeKey(e.Key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	aead, err := makeAEAD(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nonce := e.getNonce(aead.NonceSize())

	return aead.Seal(nonce, nonce, data, nil), nil
}

func (e *Encrypter) getNonce(size int) []byte {
	return generateRandom(size)
}

type Decrypter struct {
	keySetter
}

func NewDecrypter() *Decrypter {
	return &Decrypter{}
}

func (d *Decrypter) Decrypt(data []byte) ([]byte, error) {
	const op = "Decrypter.Decrypt"

	key, err := makeKey(d.Key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	aead, err := makeAEAD(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nonce, payload := data[:aead.NonceSize()], data[aead.NonceSize():]

	decData, err := aead.Open(nil, nonce, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return decData, nil
}

type keySetter struct {
	Key string
}

func (s *keySetter) SetKey(key string) {
	s.Key = key
}

func makeAEAD(key []byte) (cipher.AEAD, error) {
	const op = "cipher.makeAEAD"
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return aesGCM, nil
}

func makeKey(k string) ([]byte, error) {
	const op = "cipher.makeKey"
	key, err := pbkdf2.Key(sha256.New, k, makeSalt(k), iter, keySize)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return key, nil
}

func makeSalt(k string) []byte {
	h := sha256.New()
	h.Write([]byte(k))
	return h.Sum(nil)
}

func generateRandom(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

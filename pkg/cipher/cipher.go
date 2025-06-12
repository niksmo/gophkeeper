package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/niksmo/gophkeeper/pkg/logger"
)

const (
	keySize = 32
	iter    = 4096
)

type Encrypter interface {
	Encrypt([]byte) []byte
}

type Decrypter interface {
	Decrypt([]byte) ([]byte, error)
}

type encrypter struct {
	k string
	l logger.Logger
}

func NewEncrypter(logger logger.Logger, key string) Encrypter {
	return &encrypter{key, logger}
}

func (e *encrypter) Encrypt(data []byte) []byte {
	const op = "encrypter.Encrypt"
	log := e.l.With().Str("op", op).Logger()

	key, err := makeKey(e.k)
	if err != nil {
		log.Fatal().Err(err)
	}

	aead, err := makeAEAD(key)
	if err != nil {
		log.Fatal().Err(err)
	}

	nonce := e.getNonce(aead.NonceSize())

	return aead.Seal(nonce, nonce, data, nil)
}

type decrypter struct {
	k string
	l logger.Logger
}

func NewDecrypter(logger logger.Logger, key string) Decrypter {
	return &decrypter{key, logger}
}

func (d *decrypter) Decrypt(data []byte) ([]byte, error) {
	const op = "decrypter.Decrypt"
	log := d.l.With().Str("op", op).Logger()

	key, err := makeKey(d.k)
	if err != nil {
		log.Fatal().Err(err)
	}

	aead, err := makeAEAD(key)
	if err != nil {
		log.Fatal().Err(err)
	}

	nonce, payload := data[:aead.NonceSize()], data[aead.NonceSize():]

	decData, err := aead.Open(nil, nonce, payload, nil)
	if err != nil {
		log.Debug().Err(err).Msg("failed to decrypt")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return decData, nil
}

func (e *encrypter) getNonce(size int) []byte {
	return generateRandom(size)
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

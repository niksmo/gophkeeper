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

type Encrypter struct {
	l logger.Logger
	k string
}

func NewEncrypter(l logger.Logger) *Encrypter {
	return &Encrypter{l: l}
}

func (e *Encrypter) SetKey(k string) {
	e.k = k
}

func (e *Encrypter) Encrypt(data []byte) []byte {
	const op = "encrypter.Encrypt"
	log := e.l.With().Str("op", op).Logger()

	key, err := makeKey(e.k)
	if err != nil {
		log.Debug().Err(err)
	}

	aead, err := makeAEAD(key)
	if err != nil {
		log.Debug().Err(err)
	}

	nonce := e.getNonce(aead.NonceSize())

	return aead.Seal(nonce, nonce, data, nil)
}

func (e *Encrypter) getNonce(size int) []byte {
	return generateRandom(size)
}

type Decrypter struct {
	l logger.Logger
	k string
}

func NewDecrypter(l logger.Logger) *Decrypter {
	return &Decrypter{l: l}
}

func (d *Decrypter) SetKey(k string) {
	d.k = k
}

func (d *Decrypter) Decrypt(data []byte) ([]byte, error) {
	const op = "decrypter.Decrypt"
	log := d.l.With().Str("op", op).Logger()

	key, err := makeKey(d.k)
	if err != nil {
		log.Debug().Err(err)
	}

	aead, err := makeAEAD(key)
	if err != nil {
		log.Debug().Err(err)
	}

	nonce, payload := data[:aead.NonceSize()], data[aead.NonceSize():]

	decData, err := aead.Open(nil, nonce, payload, nil)
	if err != nil {
		log.Debug().Err(err).Msg("failed to decrypt")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return decData, nil
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

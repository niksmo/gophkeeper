package hasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTooLongSrc = errors.New("source length is too long, exceeds 72 bytes")
)

type CryptoHasher struct {
	cost int
}

func NewCryptoHasher(cost int) *CryptoHasher {
	return &CryptoHasher{cost}
}

func (h *CryptoHasher) Generate(src []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(src, h.cost)
	if err != nil && errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return nil, ErrTooLongSrc
	}
	return hash, err
}

func (h *CryptoHasher) Compare(hash, src []byte) error {
	return bcrypt.CompareHashAndPassword(hash, src)
}

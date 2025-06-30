package tokenservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var signingMethod = jwt.SigningMethodHS256

type UserTokenProvider struct {
	logger   logger.Logger
	sercret  []byte
	tokenTTL time.Duration
}

func NewUsersTokenProvider(
	logger logger.Logger, secret []byte, tokenTTL time.Duration,
) UserTokenProvider {
	return UserTokenProvider{logger, secret, tokenTTL}
}

func (tp UserTokenProvider) GetTokenString(userID int) (string, error) {
	const op = "UserTokenProvider.GetTokenString"

	token := jwt.NewWithClaims(
		signingMethod,
		newClaims(userID, tp.tokenTTL),
	)
	signed, err := token.SignedString(tp.sercret)
	if err != nil {
		tp.logger.Error().Str("op", op).Msg("failed to make signed token")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return signed, nil
}

type UserTokenVerifier struct {
	logger logger.Logger
	secret []byte
}

func NewUsersTokenVerifier(
	logger logger.Logger, secret []byte,
) UserTokenVerifier {
	return UserTokenVerifier{logger, secret}
}

func (tv UserTokenVerifier) Verify(tokenStr string) (int, error) {
	const op = "UserTokenVerifier.Verify"
	var c claims
	_, err := jwt.ParseWithClaims(tokenStr, &c, tv.keyFn)
	if err != nil {
		tv.logger.Debug().Err(err).Str("op", op).Msg("failed to parse token")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return c.UserID, nil
}

func (tv UserTokenVerifier) keyFn(t *jwt.Token) (any, error) {
	if err := tv.validMethod(t.Method); err != nil {
		return nil, err
	}
	return tv.secret, nil
}

func (tv UserTokenVerifier) validMethod(method jwt.SigningMethod) error {
	if method.Alg() == signingMethod.Alg() {
		return nil
	}

	return errors.New("unexpected signing method")
}

type claims struct {
	jwt.RegisteredClaims
	UserID int `json:"uid"`
}

func newClaims(userID int, tokenTTL time.Duration) claims {
	expiresAt := jwt.NewNumericDate(time.Now().Add(tokenTTL))
	registeredClaims := jwt.RegisteredClaims{ExpiresAt: expiresAt}
	return claims{registeredClaims, userID}
}

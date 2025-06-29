package tokenservice_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/niksmo/gophkeeper/internal/server/service/tokenservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var log = logger.NewPretty("debug")
var secret = []byte("awesomeSecret")

func TestProvider(t *testing.T) {
	tokenTTL := time.Second * 5
	userID := 777
	tp := tokenservice.NewUsersTokenProvider(log, secret, tokenTTL)
	tokenStr, err := tp.GetTokenString(userID)
	require.NoError(t, err)
	assert.NotZero(t, tokenStr)
}

func TestVerifier(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		tokenTTL := time.Second * 5
		userID := 777
		tp := tokenservice.NewUsersTokenProvider(log, secret, tokenTTL)
		tokenStr, err := tp.GetTokenString(userID)
		require.NoError(t, err)

		tv := tokenservice.NewUsersTokenVerifier(log, secret)
		extractedUserID, err := tv.Verify(tokenStr)
		require.NoError(t, err)
		assert.Equal(t, userID, extractedUserID)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		tokenTTL := time.Millisecond
		userID := 777
		tp := tokenservice.NewUsersTokenProvider(log, secret, tokenTTL)
		tokenStr, err := tp.GetTokenString(userID)
		require.NoError(t, err)

		time.Sleep(time.Millisecond * 2)

		tv := tokenservice.NewUsersTokenVerifier(log, secret)
		_, err = tv.Verify(tokenStr)
		require.ErrorIs(t, err, jwt.ErrTokenExpired)
	})

	t.Run("InvalidKey", func(t *testing.T) {
		tokenStr := "give_me_the_chance"
		tv := tokenservice.NewUsersTokenVerifier(log, secret)
		_, err := tv.Verify(tokenStr)
		require.ErrorIs(t, err, jwt.ErrTokenMalformed)
	})
}

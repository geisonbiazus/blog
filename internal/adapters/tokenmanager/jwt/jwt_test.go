package jwt_test

import (
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/tokenmanager/jwt"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestTokenManager(t *testing.T) {
	t.Run("It encodes and decodes the given user ID", func(t *testing.T) {
		secret := "secret"
		manager := jwt.NewTokenManager(secret)
		userID := "user-id"

		token, err := manager.Encode(userID, 10*time.Minute)

		assert.Nil(t, err)

		decodedUserID, err := manager.Decode(token)

		assert.Nil(t, err)
		assert.Equal(t, userID, decodedUserID)
	})

	t.Run("It returns error when decoding expired token", func(t *testing.T) {
		secret := "secret"
		manager := jwt.NewTokenManager(secret)
		userID := "user-id"

		token, err := manager.Encode(userID, -10*time.Minute)

		assert.Nil(t, err)

		decodedUserID, err := manager.Decode(token)

		assert.Error(t, auth.ErrTokenExpired, err)
		assert.Equal(t, "", decodedUserID)
	})
}

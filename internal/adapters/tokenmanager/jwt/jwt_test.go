package jwt_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/tokenmanager/jwt"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestTokenManager(t *testing.T) {
	t.Run("It encodes and decodes the given user ID", func(t *testing.T) {
		secret := "secret"
		manager := jwt.NewTokenManager(secret)
		userID := "user-id"

		token, err := manager.Encode(userID)

		assert.Nil(t, err)

		decodedUserID, err := manager.Decode(token)

		assert.Nil(t, err)
		assert.Equal(t, userID, decodedUserID)
	})
}

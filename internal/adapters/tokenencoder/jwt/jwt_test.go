package jwt_test

import (
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/tokenencoder/jwt"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/stretchr/testify/assert"
)

func TestTokenEncoder(t *testing.T) {
	t.Run("It encodes and decodes the given value", func(t *testing.T) {
		secret := "secret"
		manager := jwt.NewTokenEncoder(secret)
		value := "value"

		token, err := manager.Encode(value, 10*time.Minute)

		assert.Nil(t, err)

		decodedValue, err := manager.Decode(token)

		assert.Nil(t, err)
		assert.Equal(t, value, decodedValue)
	})

	t.Run("It returns error when decoding expired token", func(t *testing.T) {
		secret := "secret"
		manager := jwt.NewTokenEncoder(secret)
		value := "value"

		token, err := manager.Encode(value, -10*time.Minute)

		assert.Nil(t, err)

		decodedValue, err := manager.Decode(token)

		assert.Error(t, auth.ErrTokenExpired, err)
		assert.Equal(t, "", decodedValue)
	})
}

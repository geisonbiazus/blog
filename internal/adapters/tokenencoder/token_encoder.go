package tokenencoder

import "github.com/geisonbiazus/blog/internal/adapters/tokenencoder/jwt"

func NewJWTTokenEncoder(secret string) *jwt.TokenEncoder {
	return jwt.NewTokenEncoder(secret)
}

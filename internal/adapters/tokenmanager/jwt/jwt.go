package jwt

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager struct {
	secret []byte
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: []byte(secret)}
}

func (m *TokenManager) Encode(userID string) (string, error) {
	claims := &userClaims{
		UserID: userID,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := t.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("error signing string on jwt.TokenManager: %w", err)
	}

	return signedToken, nil
}

func (m *TokenManager) Decode(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, ErrInvalidSigningAlgorithm
		}

		return m.secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token on jwt.TokenManager: %w", err)
	}

	claims := t.Claims.(*userClaims)

	return claims.UserID, nil
}

var ErrInvalidSigningAlgorithm = errors.New("invalid JWT signing algorithm")

type userClaims struct {
	jwt.StandardClaims
	UserID string
}

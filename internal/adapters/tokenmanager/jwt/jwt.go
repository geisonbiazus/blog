package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/geisonbiazus/blog/internal/core/auth"
)

type TokenManager struct {
	secret []byte
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: []byte(secret)}
}

func (m *TokenManager) Encode(userID string, expiresIn time.Duration) (string, error) {
	claims := &userClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt(expiresIn),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := t.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("error signing string on jwt.TokenManager: %w", err)
	}

	return signedToken, nil
}

func expiresAt(expiresIn time.Duration) int64 {
	return time.Now().Add(expiresIn).Unix()
}

func (m *TokenManager) Decode(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, ErrInvalidSigningAlgorithm
		}

		return m.secret, nil
	})

	if err != nil {
		return "", m.handleDecodingError(err)
	}

	claims := t.Claims.(*userClaims)

	return claims.UserID, nil
}

func (m *TokenManager) handleDecodingError(err error) error {
	if strings.Contains(err.Error(), "token is expired") {
		return auth.ErrTokenExpired
	}

	return fmt.Errorf("error parsing token on jwt.TokenManager: %w", err)
}

var ErrInvalidSigningAlgorithm = errors.New("invalid JWT signing algorithm")

type userClaims struct {
	jwt.StandardClaims
	UserID string
}

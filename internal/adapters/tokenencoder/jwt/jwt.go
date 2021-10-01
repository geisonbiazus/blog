package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/geisonbiazus/blog/internal/core/auth"
)

type TokenEncoder struct {
	secret []byte
}

func NewTokenEncoder(secret string) *TokenEncoder {
	return &TokenEncoder{secret: []byte(secret)}
}

func (m *TokenEncoder) Encode(value string, expiresIn time.Duration) (string, error) {
	claims := newClaims(value, expiresIn)
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := t.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("error signing string on jwt.TokenManager: %w", err)
	}

	return signedToken, nil
}

func (m *TokenEncoder) Decode(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, ErrInvalidSigningAlgorithm
		}

		return m.secret, nil
	})

	if err != nil {
		return "", m.handleDecodingError(err)
	}

	claims := t.Claims.(*jwtClaims)

	return claims.Subject, nil
}

func (m *TokenEncoder) handleDecodingError(err error) error {
	if strings.Contains(err.Error(), "token is expired") {
		return auth.ErrTokenExpired
	}

	return fmt.Errorf("error parsing token on jwt.TokenManager: %w", err)
}

var ErrInvalidSigningAlgorithm = errors.New("invalid JWT signing algorithm")

type jwtClaims struct {
	jwt.StandardClaims
}

func newClaims(sub string, expiresIn time.Duration) *jwtClaims {
	return &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   sub,
			ExpiresAt: expiresAt(expiresIn),
		},
	}
}

func expiresAt(expiresIn time.Duration) int64 {
	return time.Now().Add(expiresIn).Unix()
}

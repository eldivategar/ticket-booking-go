package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTGenerator adalah implementasi konkret dari auth.TokenGenerator
type JWTGenerator struct {
	secretKey string
}

func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{secretKey: secret}
}

func (j *JWTGenerator) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(), // 'sub' (subject) adalah standar untuk ID user
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Token berlaku 1 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

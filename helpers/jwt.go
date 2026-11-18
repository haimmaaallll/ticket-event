package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(secret, userId, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userId,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenstring, nil
}

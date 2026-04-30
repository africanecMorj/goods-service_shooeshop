package jwtpkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccess(secret []byte, userID int, role string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"role": role,
	})
	return t.SignedString(secret)
}

func Parse(secret []byte, tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
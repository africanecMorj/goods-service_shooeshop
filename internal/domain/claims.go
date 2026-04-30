package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID      int      `json:"sub"`
    Roles       []string `json:"roles"`

    jwt.RegisteredClaims
}
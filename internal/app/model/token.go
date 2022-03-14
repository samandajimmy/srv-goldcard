package model

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Token to store JWT token data
type Token struct {
	Name   string `json:"name"`
	Claims jwt.StandardClaims
}

// AccountToken to store account token data
type AccountToken struct {
	ID        int64         `json:"id,omitempty"`
	Username  string        `json:"username"`
	Password  string        `json:"password,omitempty"`
	Token     string        `json:"token"`
	ExpireAt  *time.Time    `json:"expireAt"`
	Status    *int64        `json:"status,omitempty"`
	CreatedAt *time.Time    `json:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at"`
	ExpiresAt time.Duration `json:"expiresAt" pg:"-"`
}

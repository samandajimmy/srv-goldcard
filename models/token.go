package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Token to store JWT token data
type Token struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// AccountToken to store account token data
type AccountToken struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Token     string     `json:"token"`
	ExpireAt  *time.Time `json:"expireAt"`
	Status    *int64     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

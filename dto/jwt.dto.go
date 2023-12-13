package dto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId  int    `json:"userId,omitempty"`
	Session string `json:"session"`
	Role    string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

type JWTToken struct {
	TokenString string    `json:"tokenString,omitempty"`
	ExpiredAt   time.Time `json:"expiredAt,omitempty"`
}

func (d JWTToken) GetBearer() string {
	return fmt.Sprintf("Bearer %s", d.TokenString)
}

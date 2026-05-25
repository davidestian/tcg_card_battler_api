package auth_dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
	AccountID string `json:"accountID"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

var JWTAccessSecret = []byte("accessToken")
var JWTRefreshSecret = []byte("refreshToken")

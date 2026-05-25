package middleware

import (
	"fmt"
	"net/http"
	"strings"
	web "tcg_card_battler/web-api/internal"
	auth_dto "tcg_card_battler/web-api/internal/dto/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, web.APIResponse{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &auth_dto.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
			}
			return auth_dto.JWTAccessSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, web.APIResponse{
				Success: false,
				Message: "Invalid token",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*auth_dto.AppClaims); ok && token.Valid {
			// Set the username in the gin context for subsequent handlers
			c.Set("accountID", claims.AccountID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, web.APIResponse{
				Success: false,
				Message: "Invalid token claims",
			})
			c.Abort()
		}
	}
}

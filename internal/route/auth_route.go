package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up all routes for the /users endpoint
func RegisterAuthRoutes(rg *gin.RouterGroup, h *handler.AuthHandler) {
	routes := rg.Group("/auth")
	{
		routes.POST("/login", h.Login)
		routes.POST("/refresh", h.Refresh)
	}
}

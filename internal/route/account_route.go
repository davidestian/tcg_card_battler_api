package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up all routes for the /users endpoint
func RegisterAccountRoutes(rg *gin.RouterGroup, h handler.AccountHandler) {
	routes := rg.Group("/account")
	{
		routes.GET("", h.GetAccount)
		routes.PUT("/gold", h.PutAccountGold)
	}
}

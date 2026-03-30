package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up all routes for the /users endpoint
func RegisterUnitRoutes(rg *gin.RouterGroup, h *handler.UnitHandler) {
	routes := rg.Group("/unit")
	{
		routes.GET("", h.GetUnitByCode)
		routes.GET("/next-level-path", h.GetUnitNextLevelPath)
	}
}

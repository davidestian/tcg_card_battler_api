package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up all routes for the /users endpoint
func RegisterInventoryRoutes(rg *gin.RouterGroup, h *handler.InventoryHandler) {
	routes := rg.Group("/inventory")
	{
		routes.GET("/unit", h.GetPlayerUnits)
		routes.GET("/unit/detail", h.GetInventoryPlayerUnitDetailByCode)
		routes.GET("/unit/create", h.GetEligibleUnitsToCreate)
		routes.GET("/card", h.GetAllPlayerCard)
		routes.GET("/card/unit", h.GetPlayerUnitCardByUnitCode)
		routes.POST("/unit/level-up", h.PostPlayerUnitLevelUp)
		routes.GET("/unit/prev-level", h.GetPlayerUnitPrevLevel)
		routes.POST("/unit/level/change-image", h.PostPlayerUnitLevelChangeImage)
		routes.POST("/player-unit/upgrade", h.PostPlayerUnitUpgrade)
		routes.POST("/player-unit/create", h.PostCreatePlayerUnit)
	}
}

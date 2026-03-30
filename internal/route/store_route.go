package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up all routes for the /users endpoint
func RegisterStoreRoutes(rg *gin.RouterGroup, h *handler.StoreHandler) {
	routes := rg.Group("/store")
	{
		routes.GET("/booster", h.GetAllBooster)
		routes.GET("/booster/card", h.GetAllBoosterCard)
		routes.GET("/booster/rarity-rate", h.GetBoosterRarityRate)
		routes.POST("/booster/buy-pack", h.PostBuyBoosterPack)
	}
}

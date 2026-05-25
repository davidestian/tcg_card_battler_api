package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterGeneralRoutes(rg *gin.RouterGroup, h handler.GeneralHandler) {
	routes := rg.Group("/general")
	{
		routes.GET("/elements", h.GetElements)
		routes.GET("/origins", h.GetOrigins)
	}
}

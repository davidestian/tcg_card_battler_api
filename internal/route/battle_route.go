package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up all routes for the /users endpoint
func RegisterBattleRoutes(rg *gin.RouterGroup, bh handler.BattleHandler) {
	routes := rg.Group("/battle")
	{
		routes.GET("/unit-random", bh.GetRandomEnemyBattleUnits)
		routes.GET("/player-team", bh.GetPlayerTeamUnits)
	}
}

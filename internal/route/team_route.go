package route

import (
	"tcg_card_battler/web-api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up all routes for the /users endpoint
func RegisterTeamRoutes(rg *gin.RouterGroup, th handler.TeamHandler) {
	routes := rg.Group("/team")
	{
		routes.GET("/list", th.GetPlayerTeamList)
		routes.GET("", th.GetPlayerTeam)
		routes.GET("/active/id", th.GetActivePlayerTeamID)
		routes.POST("", th.PostPlayerTeam)
		routes.PUT("/active", th.PutActivePlayerTeam)
		routes.DELETE("", th.DeletePlayerTeam)
	}
}

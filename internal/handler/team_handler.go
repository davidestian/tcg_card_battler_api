package handler

import (
	"net/http"
	"strconv"
	"strings"
	web "tcg_card_battler/web-api/internal"
	team_dto "tcg_card_battler/web-api/internal/dto/team"
	"tcg_card_battler/web-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TeamHandler interface {
	GetPlayerTeamList(c *gin.Context)
	GetPlayerTeam(c *gin.Context)
	GetActivePlayerTeamID(c *gin.Context)
	PostPlayerTeam(c *gin.Context)
	PutActivePlayerTeam(c *gin.Context)
	DeletePlayerTeam(c *gin.Context)
}

type TeamHandlerImpl struct {
	teamService service.TeamService
}

func NewTeamHandler(tsrv service.TeamService) TeamHandler {
	return &TeamHandlerImpl{teamService: tsrv}
}

func (h *TeamHandlerImpl) GetPlayerTeamList(c *gin.Context) {
	ctx := c.Request.Context()

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil || limit <= 0 {
		web.Error(c, http.StatusBadRequest, "limit must be a positive number")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil || page <= 0 {
		web.Error(c, http.StatusBadRequest, "page must be a positive number")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.teamService.GetPlayerTeam(ctx, accountID.(string), limit, page)
	if err != nil {
		web.Error(c, http.StatusOK, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *TeamHandlerImpl) GetPlayerTeam(c *gin.Context) {
	ctx := c.Request.Context()

	teamID := c.DefaultQuery("teamID", "")
	if teamID == "" {
		web.Error(c, http.StatusOK, "teamID cannot be empty")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.teamService.GetPlayerTeamByTeamID(ctx, accountID.(string), teamID)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", result)
}

func (h *TeamHandlerImpl) GetActivePlayerTeamID(c *gin.Context) {
	ctx := c.Request.Context()

	accountID, _ := c.Get("accountID")
	result, err := h.teamService.GetActivePlayerTeamID(ctx, accountID.(string))
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", result)
}

func (h *TeamHandlerImpl) PostPlayerTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var rq team_dto.PostPlayerTeamRQ
	err := c.ShouldBindJSON(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	err = h.teamService.PostPlayerTeam(ctx, accountID.(string), rq)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", nil)
}

func (h *TeamHandlerImpl) PutActivePlayerTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var rq team_dto.PutActivePlayerTeamRQ
	err := c.ShouldBindJSON(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}
	if rq.TeamID == "" {
		web.Error(c, http.StatusOK, "teamID cannot be empty")
		return
	}

	accountID, _ := c.Get("accountID")
	err = h.teamService.PutActivePlayerTeam(ctx, accountID.(string), rq.TeamID)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", nil)
}

func (h *TeamHandlerImpl) DeletePlayerTeam(c *gin.Context) {
	ctx := c.Request.Context()

	playerTeamID := c.Query("playerTeamID")
	if strings.TrimSpace(playerTeamID) == "" {
		web.Error(c, http.StatusBadRequest, "playerTeamID cannot be empty")
		return
	}

	accountID, _ := c.Get("accountID")
	err := h.teamService.DeletePlayerTeam(ctx, accountID.(string), playerTeamID)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", nil)
}

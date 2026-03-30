package handler

import (
	"net/http"
	"strconv"
	"strings"
	web "tcg_card_battler/web-api/internal"
	"tcg_card_battler/web-api/internal/service"

	"github.com/gin-gonic/gin"
)

type BattleHandler interface {
	GetRandomEnemyBattleUnits(c *gin.Context)
	GetPlayerTeamUnits(c *gin.Context)
}

type battleHandlerImpl struct {
	battleService service.BattleService
}

func NewBattleHandler(bs service.BattleService) BattleHandler {
	return &battleHandlerImpl{battleService: bs}
}

func (s *battleHandlerImpl) GetRandomEnemyBattleUnits(c *gin.Context) {
	levelString := c.QueryArray("levels")
	if len(levelString) == 0 {
		web.Error(c, http.StatusBadRequest, "levels cannot be empty")
		return
	}

	levels, err := stringsToInts(levelString)
	if err != nil {
		web.Error(c, http.StatusBadRequest, "levels must be number")
		return
	}

	evoLevelString := c.QueryArray("evoLevels")
	if len(evoLevelString) == 0 {
		web.Error(c, http.StatusBadRequest, "evoLevels cannot be empty")
		return
	}

	evoLevels, err := stringsToInts(evoLevelString)
	if err != nil {
		web.Error(c, http.StatusBadRequest, "levels must be number")
		return
	}

	res, err := s.battleService.GetRandomEnemyBattleUnits(c, levels, evoLevels)
	if err != nil {
		web.Error(c, http.StatusOK, "Failed to Generate")
		return
	}

	web.Success(c, "", res)
}

func (s *battleHandlerImpl) GetPlayerTeamUnits(c *gin.Context) {
	playerTeamID := c.Query("playerTeamID")
	if strings.TrimSpace(playerTeamID) == "" {
		web.Error(c, http.StatusBadRequest, "playerTeamID not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	res, err := s.battleService.GetPlayerTeamUnits(c, accountID.(string), playerTeamID)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	web.Success(c, "", res)
}

func stringsToInts(strings []string) ([]int, error) {
	ints := make([]int, 0, len(strings))
	for _, s := range strings {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		ints = append(ints, val)
	}
	return ints, nil
}

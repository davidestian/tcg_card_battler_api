package handler

import (
	"net/http"
	"strings"
	web "tcg_card_battler/web-api/internal"
	"tcg_card_battler/web-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	unitService service.UnitService
}

func NewUnitHandler(svc service.UnitService) *UnitHandler {
	return &UnitHandler{unitService: svc}
}

func (h *UnitHandler) GetUnitByCode(c *gin.Context) {
	ctx := c.Request.Context()

	unitCode := c.Query("unitCode")
	if strings.TrimSpace(unitCode) == "" {
		web.Error(c, http.StatusBadRequest, "unit code not valid")
		return
	}

	result, err := h.unitService.GetUnitByCode(ctx, unitCode)
	if err != nil {
		web.Error(c, http.StatusOK, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *UnitHandler) GetUnitNextLevelPath(c *gin.Context) {
	ctx := c.Request.Context()
	unitCode := c.Query("unitCode")
	if strings.TrimSpace(unitCode) == "" {
		web.Error(c, http.StatusBadRequest, "unit code not valid")
		return
	}

	result, err := h.unitService.GetUnitLevelPathByCode(ctx, unitCode)
	if err != nil {
		web.Error(c, http.StatusOK, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

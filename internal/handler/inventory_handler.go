package handler

import (
	"net/http"
	"strconv"
	"strings"
	web "tcg_card_battler/web-api/internal"
	inv_dto "tcg_card_battler/web-api/internal/dto/inventory"
	"tcg_card_battler/web-api/internal/service"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: svc}
}

func (h *InventoryHandler) GetPlayerUnits(c *gin.Context) {
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
	result, err := h.inventoryService.GetPlayerUnits(ctx, accountID.(string), limit, page)
	if err != nil {
		web.Error(c, http.StatusInternalServerError, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *InventoryHandler) GetInventoryPlayerUnitDetailByCode(c *gin.Context) {
	ctx := c.Request.Context()
	playerUnitID := c.Query("playerUnitID")
	if strings.TrimSpace(playerUnitID) == "" {
		web.Error(c, http.StatusBadRequest, "unit id cannot be empty")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.inventoryService.GetPlayerUnitDetailByCode(ctx, accountID.(string), playerUnitID)
	if err != nil {
		web.Error(c, http.StatusInternalServerError, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *InventoryHandler) GetAllPlayerCard(c *gin.Context) {
	ctx := c.Request.Context()

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil || limit <= 0 {
		web.Error(c, http.StatusBadRequest, "limit must be a positive number")
		return
	}

	price, err := strconv.Atoi(c.DefaultQuery("price", "-1"))
	if err != nil || price < 0 {
		web.Error(c, http.StatusBadRequest, "price must be a positive number")
		return
	}

	imageTypeNumber, err := strconv.Atoi(c.DefaultQuery("imageTypeNumber", "-1"))
	if err != nil || imageTypeNumber < 0 {
		web.Error(c, http.StatusBadRequest, "imageTypeNumber must be a positive number")
		return
	}

	pageNumber, err := strconv.Atoi(c.DefaultQuery("pageNumber", "-1"))
	if err != nil || pageNumber <= 0 {
		web.Error(c, http.StatusBadRequest, "pageNumber must be a positive number")
		return
	}

	isPrev, err := strconv.ParseBool(c.DefaultQuery("isPrev", "false"))
	if err != nil {
		web.Error(c, http.StatusBadRequest, "isPrev must be boolean")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.inventoryService.GetAllPlayerCards(ctx, accountID.(string), limit, price, c.Query("code"), imageTypeNumber, pageNumber, isPrev)
	if err != nil {
		web.Error(c, http.StatusInternalServerError, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *InventoryHandler) GetPlayerUnitCardByUnitCode(c *gin.Context) {
	ctx := c.Request.Context()
	unitCode := c.Query("unitCode")
	if strings.TrimSpace(unitCode) == "" {
		web.Error(c, http.StatusBadRequest, "unit code not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.inventoryService.GetPlayerUnitCardByUnitCode(ctx, accountID.(string), unitCode)
	if err != nil {
		web.Error(c, http.StatusInternalServerError, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *InventoryHandler) PostPlayerUnitLevelUp(c *gin.Context) {
	ctx := c.Request.Context()
	var rq inv_dto.PostPlayerUnitLevelUpRQ

	err := c.ShouldBindJSON(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	err = h.inventoryService.PostPlayerUnitLevelUp(ctx, accountID.(string), rq)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", nil)
}

func (h *InventoryHandler) GetPlayerUnitPrevLevel(c *gin.Context) {
	ctx := c.Request.Context()
	playerUnitID := c.Query("playerUnitID")

	if strings.TrimSpace(playerUnitID) == "" {
		web.Error(c, http.StatusBadRequest, "player unit id not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	results, err := h.inventoryService.GetPlayerUnitPrevLevel(ctx, accountID.(string), playerUnitID)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", results)
}

func (h *InventoryHandler) PostPlayerUnitLevelChangeImage(c *gin.Context) {
	ctx := c.Request.Context()
	var rq inv_dto.PlayerUnitLevelChangeImageRQ

	err := c.ShouldBindJSON(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	err = h.inventoryService.ChangePlayerUnitImage(ctx, accountID.(string), rq)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", nil)
}

func (h *InventoryHandler) PostPlayerUnitUpgrade(c *gin.Context) {
	ctx := c.Request.Context()
	var rq inv_dto.PostPlayerUnitUpgradeRQ

	err := c.ShouldBindJSON(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	err = h.inventoryService.PostPlayerUnitUpgrade(ctx, accountID.(string), rq)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", nil)
}

func (h *InventoryHandler) GetEligibleUnitsToCreate(c *gin.Context) {
	ctx := c.Request.Context()
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil || limit <= 0 {
		web.Error(c, http.StatusBadRequest, "limit must be a positive number")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "-1"))
	if err != nil || page <= 0 {
		web.Error(c, http.StatusBadRequest, "pageNumber must be a positive number")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.inventoryService.GetEligibleUnitsToCreate(ctx, accountID.(string), limit, page)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", result)
}

func (h *InventoryHandler) PostCreatePlayerUnit(c *gin.Context) {
	ctx := c.Request.Context()
	var rq inv_dto.PostCreatePlayerUnitRQ

	err := c.ShouldBindJSON(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	err = h.inventoryService.PostCreatePlayerUnit(ctx, accountID.(string), rq)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "", nil)
}

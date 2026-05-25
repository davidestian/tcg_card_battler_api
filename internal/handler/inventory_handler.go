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

type InventoryHandler interface {
	GetPlayerUnits(c *gin.Context)
	GetInventoryPlayerUnitDetailByCode(c *gin.Context)
	GetAllPlayerCard(c *gin.Context)
	GetPlayerUnitCardByUnitCode(c *gin.Context)
	PostPlayerUnitLevelUp(c *gin.Context)
	GetPlayerUnitPrevLevel(c *gin.Context)
	PostPlayerUnitLevelChangeImage(c *gin.Context)
	PostPlayerUnitUpgrade(c *gin.Context)
	GetEligibleUnitsToCreate(c *gin.Context)
	PostCreatePlayerUnit(c *gin.Context)
	GetPlayerUnitPriceByID(c *gin.Context)
	PostSellPlayerUnit(c *gin.Context)
}

type InventoryHandlerImpl struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(svc service.InventoryService) InventoryHandler {
	return &InventoryHandlerImpl{inventoryService: svc}
}

func (h *InventoryHandlerImpl) GetPlayerUnits(c *gin.Context) {
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

	level, _ := strconv.Atoi(c.DefaultQuery("level", "0"))
	lastUnitLevel, _ := strconv.Atoi(c.DefaultQuery("lastUnitLevel", "0"))
	element1, _ := strconv.Atoi(c.DefaultQuery("element1", "0"))
	element2, _ := strconv.Atoi(c.DefaultQuery("element2", "0"))
	sort, _ := strconv.Atoi(c.DefaultQuery("sort", "0"))
	name := c.DefaultQuery("name", "")
	origin := c.DefaultQuery("origin", "")
	accountID, _ := c.Get("accountID")

	result, err := h.inventoryService.GetPlayerUnits(ctx, accountID.(string), limit, page,
		element1, element2, level, lastUnitLevel, sort, origin, name)
	if err != nil {
		web.Error(c, http.StatusInternalServerError, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *InventoryHandlerImpl) GetInventoryPlayerUnitDetailByCode(c *gin.Context) {
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

func (h *InventoryHandlerImpl) GetAllPlayerCard(c *gin.Context) {
	ctx := c.Request.Context()

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil || limit <= 0 {
		web.Error(c, http.StatusBadRequest, "limit must be a positive number")
		return
	}

	pageNumber, err := strconv.Atoi(c.DefaultQuery("pageNumber", "-1"))
	if err != nil || pageNumber <= 0 {
		web.Error(c, http.StatusBadRequest, "pageNumber must be a positive number")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.inventoryService.GetAllPlayerCards(ctx, accountID.(string), limit, pageNumber)
	if err != nil {
		web.Error(c, http.StatusInternalServerError, "failed to gets")
		return
	}
	web.Success(c, "", result)
}

func (h *InventoryHandlerImpl) GetPlayerUnitCardByUnitCode(c *gin.Context) {
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

func (h *InventoryHandlerImpl) PostPlayerUnitLevelUp(c *gin.Context) {
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

func (h *InventoryHandlerImpl) GetPlayerUnitPrevLevel(c *gin.Context) {
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

func (h *InventoryHandlerImpl) PostPlayerUnitLevelChangeImage(c *gin.Context) {
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

func (h *InventoryHandlerImpl) PostPlayerUnitUpgrade(c *gin.Context) {
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

func (h *InventoryHandlerImpl) GetEligibleUnitsToCreate(c *gin.Context) {
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

func (h *InventoryHandlerImpl) PostCreatePlayerUnit(c *gin.Context) {
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

func (h *InventoryHandlerImpl) GetPlayerUnitPriceByID(c *gin.Context) {
	ctx := c.Request.Context()
	playerUnitID := c.Query("playerUnitID")
	if strings.TrimSpace(playerUnitID) == "" {
		web.Error(c, http.StatusOK, "playerUnitID required")
		return
	}

	gold, err := h.inventoryService.GetPlayerUnitPriceByID(ctx, playerUnitID)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "Success", gold)
}

func (h *InventoryHandlerImpl) PostSellPlayerUnit(c *gin.Context) {
	ctx := c.Request.Context()
	playerUnitID := c.Query("playerUnitID")
	if strings.TrimSpace(playerUnitID) == "" {
		web.Error(c, http.StatusOK, "playerUnitID required")
		return
	}

	accountID, _ := c.Get("accountID")
	err := h.inventoryService.SellPlayerUnit(ctx, accountID.(string), playerUnitID)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}
	web.Success(c, "Success", nil)
}

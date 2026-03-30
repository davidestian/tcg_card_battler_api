package handler

import (
	"net/http"
	"strings"
	web "tcg_card_battler/web-api/internal"
	store_dto "tcg_card_battler/web-api/internal/dto/store"
	"tcg_card_battler/web-api/internal/service"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	boosterService service.BoosterService
	storeService   service.StoreService
}

func NewStoreHandler(bsrv service.BoosterService, ssrv service.StoreService) *StoreHandler {
	return &StoreHandler{boosterService: bsrv, storeService: ssrv}
}

func (h *StoreHandler) GetAllBooster(c *gin.Context) {
	ctx := c.Request.Context()

	result, err := h.boosterService.GetAllBooster(ctx)
	if err != nil {
		web.Error(c, http.StatusOK, "failed to gets")
		return
	}

	web.Success(c, "", result)
}

func (h *StoreHandler) GetAllBoosterCard(c *gin.Context) {
	ctx := c.Request.Context()

	boosterCode := c.Query("boosterCode")
	if strings.TrimSpace(boosterCode) == "" {
		web.Error(c, http.StatusBadRequest, "booster code cannot be empty")
		return
	}

	result, err := h.boosterService.GetAllBoosterCard(ctx, boosterCode)
	if err != nil {
		web.Error(c, http.StatusOK, "failed to gets")
		return
	}

	web.Success(c, "", result)
}

func (h *StoreHandler) GetBoosterRarityRate(c *gin.Context) {
	ctx := c.Request.Context()

	boosterCode := c.Query("boosterCode")
	if strings.TrimSpace(boosterCode) == "" {
		web.Error(c, http.StatusBadRequest, "booster code cannot be empty")
		return
	}

	result, err := h.boosterService.GetBoosterRarityRate(ctx, boosterCode)
	if err != nil {
		web.Error(c, http.StatusOK, "failed to gets")
		return
	}

	web.Success(c, "", result)
}

func (h *StoreHandler) PostBuyBoosterPack(c *gin.Context) {
	ctx := c.Request.Context()
	var rq store_dto.PostBuyBoosterPackRQ

	err := c.ShouldBindJSON(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	result, err := h.storeService.PostBuyBoosterPack(ctx, accountID.(string), rq)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	web.Success(c, "", result)
}

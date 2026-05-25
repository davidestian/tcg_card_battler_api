package handler

import (
	"net/http"
	web "tcg_card_battler/web-api/internal"
	"tcg_card_battler/web-api/internal/service"

	"github.com/gin-gonic/gin"
)

type GeneralHandler interface {
	GetElements(c *gin.Context)
	GetOrigins(c *gin.Context)
}

type generalHandlerImpl struct {
	generalService service.GeneralService
}

func NewGeneralHandler(g service.GeneralService) GeneralHandler {
	return &generalHandlerImpl{generalService: g}
}

func (h *generalHandlerImpl) GetElements(c *gin.Context) {
	elements, err := h.generalService.GetElements(c)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	web.Success(c, "Success", elements)
}

func (h *generalHandlerImpl) GetOrigins(c *gin.Context) {
	origins, err := h.generalService.GetOrigins(c)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	web.Success(c, "Success", origins)
}

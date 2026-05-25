package handler

import (
	"net/http"
	web "tcg_card_battler/web-api/internal"
	account_dto "tcg_card_battler/web-api/internal/dto/account"
	"tcg_card_battler/web-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	GetAccount(c *gin.Context)
	PutAccountGold(c *gin.Context)
}

type AccountHandlerImpl struct {
	accountService service.AccountService
}

func NewAccountHandler(svc service.AccountService) AccountHandler {
	return &AccountHandlerImpl{accountService: svc}
}

func (h *AccountHandlerImpl) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()
	accountID, _ := c.Get("accountID")
	data, err := h.accountService.GetAccountByID(ctx, accountID.(string))
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	web.Success(c, "Success", data)
}

func (h *AccountHandlerImpl) PutAccountGold(c *gin.Context) {
	ctx := c.Request.Context()
	var rq account_dto.PutAccountGoldRQ
	err := c.ShouldBind(&rq)
	if err != nil {
		web.Error(c, http.StatusOK, "request not valid")
		return
	}

	accountID, _ := c.Get("accountID")
	err = h.accountService.UpdateAccountGold(ctx, accountID.(string), rq.Gold)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	web.Success(c, "Success", nil)
}

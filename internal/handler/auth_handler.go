package handler

import (
	"fmt"
	"net/http"
	web "tcg_card_battler/web-api/internal"
	account_dto "tcg_card_battler/web-api/internal/dto/account"
	auth_dto "tcg_card_battler/web-api/internal/dto/auth"
	"tcg_card_battler/web-api/internal/service"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Refresh(c *gin.Context)
	Create(c *gin.Context)
	ForgotPassword(c *gin.Context)
}

type AuthHandlerImpl struct {
	accountService service.AccountService
}

func NewAuthHandler(svc service.AccountService) AuthHandler {
	return &AuthHandlerImpl{accountService: svc}
}

func generateJWT(m account_dto.AccountDetailRS) (string, string, error) {
	claims := auth_dto.AppClaims{
		AccountID: m.AccountID,
		Email:     m.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Subject:   m.AccountID,
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(auth_dto.JWTAccessSecret)
	if err != nil {
		return "", "", err
	}

	rtClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		Subject:   m.AccountID,
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString(auth_dto.JWTRefreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (h *AuthHandlerImpl) Login(c *gin.Context) {
	var input auth_dto.LoginRQ

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	result, err := h.accountService.GetAccountByEmail(c.Request.Context(), input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	} else if result == nil || result.AccountID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Account not found"})
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, result.PasswordHash)
	if err != nil || !match {
		c.JSON(http.StatusBadRequest, gin.H{"message": "login failed"})
		return
	}

	accessToken, refreshToken, err := generateJWT(account_dto.AccountDetailRS{
		AccountID: result.AccountID.String(),
		Email:     result.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	web.Success(c, "", auth_dto.LoginRS{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandlerImpl) Refresh(c *gin.Context) {
	var input auth_dto.RefreshRQ

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	// Parse and validate the Refresh Token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return auth_dto.JWTRefreshSecret, nil
	})

	if err != nil || !token.Valid {
		fmt.Println("invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	accountID := claims["sub"].(string)

	result, _ := h.accountService.GetAccountByID(c.Request.Context(), accountID)
	newAT, newRT, _ := generateJWT(*result)

	web.Success(c, "", auth_dto.LoginRS{
		AccessToken:  newAT,
		RefreshToken: newRT,
	})
}

func (h *AuthHandlerImpl) Create(c *gin.Context) {
	var input auth_dto.CreateAccountRQ

	err := c.ShouldBindJSON(&input)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	err = h.accountService.CreateNewAccount(c, input.Email, input.Username, input.Password)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	web.Success(c, "Success", nil)
}

func (h *AuthHandlerImpl) ForgotPassword(c *gin.Context) {
	var input auth_dto.ForgotPasswordRQ

	err := c.ShouldBindJSON(&input)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	err = h.accountService.ForgotPassword(c, input.Email, input.Password)
	if err != nil {
		web.Error(c, http.StatusOK, err.Error())
		return
	}

	web.Success(c, "Success", nil)
}

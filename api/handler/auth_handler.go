package handler

import (
	"manga/api/helper"
	"manga/config"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"manga/internal/usecase"

	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	LU       usecase.ILoginUsecase
	Cfg      *config.Config
	OidcCfg  *oauth2.Config
	Provider *oidc.Provider
}

func NewAuthHandler(LU usecase.ILoginUsecase, cfg *config.Config, OidcCfg *oauth2.Config, Provider *oidc.Provider) *AuthHandler {
	return &AuthHandler{
		LU:       LU,
		Cfg:      cfg,
		OidcCfg:  OidcCfg,
		Provider: Provider}
}

func (lc *AuthHandler) Login(c *gin.Context) {
	var request dtos.AuthRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	verifier := lc.Provider.Verifier(&oidc.Config{ClientID: lc.OidcCfg.ClientID})

	idToken, err := verifier.Verify(c, request.IDToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err))
		return
	}
	var claims models.GoogleClaims
	if err := idToken.Claims(&claims); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err))
		return
	}

	user, str, err := lc.LU.Login(c, claims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err))
		return
	}
	if str == "created" {
		c.JSON(http.StatusCreated, user)
	} else if str == "exist" {
		c.JSON(http.StatusOK, user)
	}
}

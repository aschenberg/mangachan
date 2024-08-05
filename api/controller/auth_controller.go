package controller

import (
	"manga/config"

	"manga/api"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"manga/pkg"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

type AuthController struct {
	LoginUsecase domain.LoginUsecase
	Conf         *config.Config
	Oidc         *oauth2.Config
	Provider     *oidc.Provider
}

func (lc *AuthController) Login(c *gin.Context) {
	var request dtos.AuthRequest

	err := c.ShouldBind(&request)
	if err != nil {
		api.RenderErrorResponse(c, "invalid request",
			pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	// token, err := lc.Oidc.Exchange(c, request.IDToken, oauth2.AccessTypeOffline)
	// if err != nil {
	// 	api.RenderErrorResponse(c, "Failed to exchange token:",
	// 		pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "invalid token"))
	// 	return
	// }

	// rawIDToken, ok := token.Extra("id_token").(string)
	// if !ok {
	// 	api.RenderErrorResponse(c, "No id_token in token response",
	// 		pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "No id_token in token response"))
	// 	return
	// }

	verifier := lc.Provider.Verifier(&oidc.Config{ClientID: lc.Oidc.ClientID})

	idToken, err := verifier.Verify(c, request.IDToken)
	if err != nil {
		api.RenderErrorResponse(c, "Failed to verify ID token: ",
			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "Failed to verify ID token: "+err.Error()))

		return
	}
	var claims struct {
		Iss           string `json:"iss"`
		Azp           string `json:"azp"`
		Aud           string `json:"aud"`
		Sub           string `json:"sub"`
		Athash        string `json:"at_hash"`
		Hd            string `json:"hd"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Iat           int    `json:"iat"`
		Exp           int    `json:"exp"`
	}
	if err := idToken.Claims(&claims); err != nil {
		api.RenderErrorResponse(c, "Failed to parse claims: ",
			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "Failed to parse claims: "+err.Error()))
		return
	}

	user, err := lc.LoginUsecase.GetUserBySub(c, claims.Sub)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			usr := models.User{
				ID:         primitive.NewObjectID(),
				AppID:      claims.Sub,
				Email:      claims.Email,
				Name:       claims.Name,
				Picture:    claims.Picture,
				Role:       []string{},
				GivenName:  claims.GivenName,
				FamilyName: claims.FamilyName}

			useCreate, err := lc.LoginUsecase.CreateUser(c, usr)
			if err != nil {
				api.RenderErrorResponse(c, "An error occurred: ",
					pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "An error occurred: "+err.Error()))
				return
			}
			GenerateToken(c, lc, useCreate)

		} else {
			api.RenderErrorResponse(c, "An error occurred: ",
				pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "An error occurred: "+err.Error()))
		}
		return
	}

	GenerateToken(c, lc, user)
}

func GenerateToken(c *gin.Context, lc *AuthController, usr models.User) {
	accessToken, err := lc.LoginUsecase.CreateAccessToken(usr, lc.Conf.JWT.AccessTokenSecret, lc.Conf.JWT.AccessTokenExpireHour)
	if err != nil {
		api.RenderErrorResponse(c, "Failed create access token: ",
			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "Failed create access token: "+err.Error()))
		return
	}

	refreshToken, err := lc.LoginUsecase.CreateRefreshToken(usr, lc.Conf.JWT.RefreshTokenSecret, lc.Conf.JWT.RefreshTokenExpireHour)
	if err != nil {
		api.RenderErrorResponse(c, "Failed create refresh token: ",
			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "Failed create refresh token: "+err.Error()))
		return
	}

	err = lc.LoginUsecase.UpdateRefresh(c, usr.ID.Hex(), refreshToken)
	if err != nil {
		api.RenderErrorResponse(c, "User Not Found with Refresh Token: ",
			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "User Not Found with Refresh Token: "+err.Error()))
		return
	}

	loginResponse := dtos.LoginResponse{
		ID:           usr.ID.Hex(),
		Email:        usr.Email,
		Picture:      usr.Picture,
		Role:         usr.Role,
		Name:         usr.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, loginResponse)
}

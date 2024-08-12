package usecase

import (
	"context"
	"manga/config"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"time"
)

type ILoginUsecase interface {
	Login(c context.Context, claims models.GoogleClaims) (dtos.LoginResponse, string, error)
}

type loginUsecase struct {
	UR             domain.UserRepository
	contextTimeout time.Duration
	Cfg            *config.Config
}

func NewLoginUsecase(UR domain.UserRepository, timeout time.Duration) ILoginUsecase {
	return &loginUsecase{
		UR:             UR,
		contextTimeout: timeout,
	}
}

func (lu *loginUsecase) Login(c context.Context, claims models.GoogleClaims) (dtos.LoginResponse, string, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	// isUserExist, err := lu.UR.IsExistBySubID(ctx, claims.Sub)
	// if err != nil {
	// 	print(err.Error())
	// 	return dtos.LoginResponse{}, "", fmt.Errorf("failed check user exist: %w", err)
	// }
	// print(isUserExist)
	// // if !isUserExist {
	usr := models.User{AppID: claims.Sub, Email: claims.Email, Picture: claims.Picture, GivenName: claims.GivenName, FamilyName: claims.FamilyName, Name: claims.Name}
	_, err := lu.UR.Create(ctx, usr)
	if err != nil {
		return dtos.LoginResponse{}, "", err
	}
	// refreshToken, err := tokenutil.CreateRefreshToken(
	// 	usr, lu.Cfg.JWT.RefreshTokenSecret,
	// 	lu.Cfg.JWT.RefreshTokenExpireHour)
	// if err != nil {
	// 	return dtos.LoginResponse{}, "", err
	// }
	// err = lu.UR.UpdateRefreshToken(ctx, claims.Sub, refreshToken)
	// if err != nil {
	// 	return dtos.LoginResponse{}, "", err
	// }
	// accessToken, err := tokenutil.CreateAccessToken(
	// 	usr, lu.Cfg.JWT.RefreshTokenSecret,
	// 	lu.Cfg.JWT.RefreshTokenExpireHour)
	// if err != nil {
	// 	return dtos.LoginResponse{}, "", err
	// }

	response := dtos.LoginResponse{
		AccessToken:  "",
		RefreshToken: "",
		ID:           "created.ID",
		Name:         "created.Name",
		Email:        "created.Email",
		Picture:      "created.Picture",
	}
	return response, "created", nil
	// } else {
	// 	usr, err := lu.UR.UpdateBySubID(ctx, claims)
	// 	if err != nil {
	// 		return dtos.LoginResponse{}, "", err
	// 	}
	// 	refreshToken, err := tokenutil.CreateRefreshToken(
	// 		usr, lu.Cfg.JWT.RefreshTokenSecret,
	// 		lu.Cfg.JWT.RefreshTokenExpireHour)
	// 	if err != nil {
	// 		return dtos.LoginResponse{}, "", err
	// 	}
	// 	err = lu.UR.UpdateRefreshToken(ctx, claims.Sub, refreshToken)
	// 	if err != nil {
	// 		return dtos.LoginResponse{}, "", err
	// 	}
	// 	accessToken, err := tokenutil.CreateAccessToken(
	// 		usr, lu.Cfg.JWT.RefreshTokenSecret,
	// 		lu.Cfg.JWT.RefreshTokenExpireHour)
	// 	if err != nil {
	// 		return dtos.LoginResponse{}, "", err
	// 	}
	// 	response := dtos.LoginResponse{
	// 		AccessToken:  accessToken,
	// 		RefreshToken: refreshToken,
	// 		ID:           usr.ID,
	// 		Name:         usr.Name,
	// 		Email:        usr.Email,
	// 		Picture:      usr.Picture,
	// 	}
	// 	return response, "exist", nil
	// }
}

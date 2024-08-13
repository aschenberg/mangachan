package usecase

import (
	"context"
	"manga/config"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"manga/pkg/tokenutil"
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

	usr, status, err := lu.UR.CreateOrUpdate(ctx, claims)
	if err != nil {
		return dtos.LoginResponse{}, "", err
	}
	refreshToken, err := tokenutil.CreateRefreshToken(
		usr, lu.Cfg.JWT.RefreshTokenSecret,
		lu.Cfg.JWT.RefreshTokenExpireHour)
	if err != nil {
		return dtos.LoginResponse{}, "", err
	}
	accessToken, err := tokenutil.CreateAccessToken(
		usr, lu.Cfg.JWT.RefreshTokenSecret,
		lu.Cfg.JWT.RefreshTokenExpireHour)
	if err != nil {
		return dtos.LoginResponse{}, "", err
	}

	response := dtos.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Name,
		Picture:      usr.Picture,
	}

	return response, status, nil
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

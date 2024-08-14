package usecase

import (
	"context"
	"manga/config"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"manga/internal/infra/pgsql/pgdb"
	"manga/pkg/flake"
	"manga/pkg/logging"
	"manga/pkg/tokenutil"
	"manga/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ILoginUsecase interface {
	Login(c context.Context, claims models.GoogleClaims) (dtos.LoginResponse, string, error)
}

type loginUsecase struct {
	UR             domain.UserRepository
	contextTimeout time.Duration
	Cfg            *config.Config
	Log            logging.Logger
}

func NewLoginUsecase(UR domain.UserRepository, timeout time.Duration, log logging.Logger) ILoginUsecase {
	return &loginUsecase{
		UR:             UR,
		contextTimeout: timeout,
		Log:            log,
	}
}

func (uc *loginUsecase) Login(c context.Context, claims models.GoogleClaims) (dtos.LoginResponse, string, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	createTime := time.Now().UTC().UnixMilli()
	key, err := flake.GenerateID(1, 1, 1)
	if err != nil {
		uc.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
		return dtos.LoginResponse{}, "", err
	}

	//TODO - Create Or Update User
	created := pgdb.CreateOrUpdateUserParams{
		ID: key, Appid: claims.Sub, Email: claims.Email,
		Picture: pgtype.Text{
			String: claims.Picture,
			Valid:  utils.StrIsEmpty(claims.Picture)},
		Role:     0,
		Isactive: true,
		GivenName: pgtype.Text{
			String: claims.GivenName,
			Valid:  utils.StrIsEmpty(claims.GivenName)},
		FamilyName: pgtype.Text{
			String: claims.FamilyName,
			Valid:  utils.StrIsEmpty(claims.FamilyName)},
		Name: pgtype.Text{
			String: claims.Name,
			Valid:  utils.StrIsEmpty(claims.Name)}, Givenname: pgtype.Text{
			String: claims.GivenName,
			Valid:  utils.StrIsEmpty(claims.GivenName)}, Familyname: pgtype.Text{
			String: claims.FamilyName,
			Valid:  utils.StrIsEmpty(claims.FamilyName)},
		Createdat: createTime,
		Updatedat: createTime}

	usr, err := uc.UR.CreateOrUpdate(ctx, created)
	if err != nil {

		return dtos.LoginResponse{}, "", err
	}

	//TODO - Generated Refresh Token
	user := models.User{
		UserID:     utils.Int64ToStr(usr.UserID),
		GivenName:  usr.GivenName.String,
		FamilyName: usr.FamilyName.String}

	refreshToken, err := tokenutil.CreateRefreshToken(
		user, uc.Cfg.JWT.RefreshTokenSecret,
		uc.Cfg.JWT.RefreshTokenExpireHour)
	if err != nil {
		uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return dtos.LoginResponse{}, "", err
	}
	tokenPr := pgdb.UpdateRefreshTokenParams{
		Refreshtoken: refreshToken, AppID: claims.Sub,
	}

	//TODO - Updated User Refresh Token
	err = uc.UR.UpdateRefreshToken(ctx, tokenPr)
	if err != nil {
		return dtos.LoginResponse{}, "", err
	}

	//TODO - Generated Access Token
	accessToken, err := tokenutil.CreateAccessToken(
		user, uc.Cfg.JWT.RefreshTokenSecret,
		uc.Cfg.JWT.RefreshTokenExpireHour)
	if err != nil {
		uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return dtos.LoginResponse{}, "", err
	}

	response := dtos.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ID:           utils.Int64ToStr(usr.UserID),
		Name:         usr.Name.String,
		Email:        usr.Email,
		Picture:      usr.Picture.String,
	}

	return response, usr.Operation, nil
}

package domain

import (
	"context"
	"manga/domain/models"
)

type LoginUsecase interface {
	GetUserBySub(c context.Context, sub string) (models.User, error)
	CreateAccessToken(user models.User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user models.User, secret string, expiry int) (refreshToken string, err error)
	CreateUser(c context.Context, user models.User) (models.User, error)
	UpdateRefresh(c context.Context, id string, token string) error
}

package domain

import (
	"context"
	"manga/internal/domain/models"
)

type SignupUsecase interface {
	Create(c context.Context, user models.User) error
	GetUserByEmail(c context.Context, email string) (models.User, error)
	CreateAccessToken(user models.User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user models.User, secret string, expiry int) (refreshToken string, err error)
}

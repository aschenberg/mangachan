package usecase

import (
	"context"
	"manga/domain"
	"manga/domain/models"
	"manga/pkg/tokenutil"
	"time"
)

type loginUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *loginUsecase) GetUserBySub(c context.Context, sub string) (models.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByAppId(ctx, sub)
}

func (lu *loginUsecase) CreateUser(c context.Context, user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.Create(ctx, user)
}

func (lu *loginUsecase) CreateAccessToken(user models.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (lu *loginUsecase) CreateRefreshToken(user models.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
func (lu *loginUsecase) UpdateRefresh(c context.Context, id string, token string) error {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.UpdateRefreshToken(ctx, id, token)
}

package usecase

import (
	"context"
	"manga/internal/domain"
	"manga/internal/domain/models"
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

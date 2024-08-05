package domain

import (
	"context"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
)

type UserRepository interface {
	Create(c context.Context, user models.User) (models.User, error)
	Fetch(c context.Context) ([]models.User, error)
	GetByAppId(c context.Context, appId string) (models.User, error)
	GetByID(c context.Context, id string) (models.User, error)
	FindByRefreshToken(c context.Context, token string) (models.User, error)
	UpdateRefreshToken(c context.Context, id string, token string) error
	DeleteByID(c context.Context, id string) error
}

type UserUsecase interface {
	GetUserByName(c context.Context, name string, page int64, pageSize int64, rsid string) ([]dtos.UserResponse, dtos.Pagination, error)
	DeleteByID(c context.Context, id string) error
}

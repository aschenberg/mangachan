package domain

import (
	"context"
	"manga/internal/domain/models"
)

const (
	CollectionUser = "users"
)

type UserRepository interface {
	Create(c context.Context, user models.User) (models.User, error)
	UpdateBySubID(c context.Context, claim models.GoogleClaims) (models.User, error)
	IsExistBySubID(c context.Context, subId string) (bool, error)
	UpdateRefreshToken(c context.Context, subId string, token string) error
	CreateOrUpdate(c context.Context, claim models.GoogleClaims) (models.User, string, error)
}

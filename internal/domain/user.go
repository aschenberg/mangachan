package domain

import (
	"context"
	"manga/internal/infra/pgsql/pgdb"
)

const (
	CollectionUser = "users"
)

type UserRepository interface {
	CreateOrUpdate(c context.Context, user pgdb.CreateOrUpdateUserParams) (pgdb.CreateOrUpdateUserRow, error)
	UpdateRefreshToken(c context.Context, token pgdb.UpdateRefreshTokenParams) error
}

package domain

import (
	"context"
	"manga/internal/infra/pgsql/pgdb"
)

const (
	CollectionUser = "users"
)

type UserRepository interface {
	CreateOrUpdate(c context.Context, user pgdb.CreateUserParams) (pgdb.CreateUserRow, error)
}

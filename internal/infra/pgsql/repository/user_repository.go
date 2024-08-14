package pgsql

import (
	"context"
	"manga/internal/domain"
	"manga/internal/infra/pgsql/pgdb"
	"manga/pkg/postgres"
)

type userRepository struct {
	q *pgdb.Queries
}

func NewUserRepository(pg *postgres.Postgres) domain.UserRepository {
	return &userRepository{
		q: pgdb.New(pg.Pool),
	}
}

func (ur *userRepository) CreateOrUpdate(c context.Context, user pgdb.CreateUserParams) (pgdb.CreateUserRow, error) {
	usr, err := ur.q.CreateUser(c, user)
	if err != nil {
		return pgdb.CreateUserRow{}, err
	}
	return usr, nil
}

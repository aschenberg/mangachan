package pgsql

import (
	"context"
	"manga/internal/domain"
	"manga/internal/infra/pgsql/pgdb"
	"manga/pkg"
	"manga/pkg/logging"
)

type userRepository struct {
	q   *pgdb.Queries
	Log logging.Logger
}

func NewUserRepository(pg *pkg.Postgres, Log logging.Logger) domain.IUserRepository {
	return &userRepository{
		q:   pgdb.New(pg.Pool),
		Log: Log,
	}
}

func (rp *userRepository) CreateOrUpdate(c context.Context, user pgdb.CreateOrUpdateUserParams) (pgdb.CreateOrUpdateUserRow, error) {
	usr, err := rp.q.CreateOrUpdateUser(c, user)
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.Upsert, err.Error(), nil)
		return pgdb.CreateOrUpdateUserRow{}, err
	}
	return usr, nil
}

func (rp *userRepository) UpdateRefreshToken(c context.Context, token pgdb.UpdateRefreshTokenParams) error {
	err := rp.q.UpdateRefreshToken(c, token)
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return err
	}
	return nil
}

package pgsql

import (
	"context"
	"manga/internal/domain"
	"manga/internal/infra/pgsql/pgdb"
	pkg "manga/pkg"
	"manga/pkg/logging"
	"manga/pkg/utils"
)

type mangaRepository struct {
	q   *pgdb.Queries
	Log logging.Logger
}

func NewMangaRepository(pg *pkg.Postgres, Log logging.Logger) domain.IMangaRepository {
	return &mangaRepository{
		q:   pgdb.New(pg.Pool),
		Log: Log,
	}
}

func (rp *mangaRepository) CreateCover(c context.Context, cvr pgdb.CreateMangaCoverParams) (int64, error) {
	id, err := rp.q.CreateMangaCover(c, cvr)
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return 0, err
	}
	return id, nil
}

func (rp *mangaRepository) FindbyId(c context.Context, mangaId string) (pgdb.FindMangaByIDRow, error) {
	res, err := rp.q.FindMangaByID(c, utils.StrToInt64(mangaId))
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return pgdb.FindMangaByIDRow{}, err
	}
	return res, nil
}
func (rp *mangaRepository) FindMangaGenre(c context.Context, mangaId string) ([]pgdb.FindMangaGenreRow, error) {
	res, err := rp.q.FindMangaGenre(c, utils.StrToInt64(mangaId))
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return []pgdb.FindMangaGenreRow{}, err
	}
	return res, nil
}

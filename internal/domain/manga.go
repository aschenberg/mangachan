package domain

import (
	"context"
	"manga/internal/infra/pgsql/pgdb"
)

type IMangaRepository interface {
	CreateCover(c context.Context, cvr pgdb.CreateMangaCoverParams) (int64, error)
	FindbyId(c context.Context, mangaId string) (pgdb.FindMangaByIDRow, error)
	FindMangaGenre(c context.Context, mangaId string) ([]pgdb.FindMangaGenreRow, error)
}

package pgsql

import (
	"context"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/infra/pgsql"
	"manga/internal/infra/pgsql/pgdb"
	"manga/pkg/flake"
	"manga/pkg/logging"
	"manga/pkg/postgres"
	"manga/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type uowRepository struct {
	q   *pgdb.Queries
	pgx *pgxpool.Pool
	Log logging.Logger
}

func NewUOWRepository(pg *postgres.Postgres, Log logging.Logger) domain.IUOWRepository {
	return &uowRepository{
		pgx: pg.Pool,
		Log: Log,
		q:   pgdb.New(pg.Pool),
	}
}

func (rp *uowRepository) CreateManga(c context.Context, u dtos.CreateManga) (dtos.CreatedMangaResponse, error) {
	tx, err := rp.pgx.Begin(c)
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.CreateTx, err.Error(), nil)
		return dtos.CreatedMangaResponse{}, err
	}

	qtx := rp.q.WithTx(tx)

	//TODO - Generated Snow Flake ID
	key, err := flake.GenerateID(1, 2, 1)
	if err != nil {
		rp.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
		tx.Rollback(c)
		return dtos.CreatedMangaResponse{}, err
	}
	//TODO - Generate Time
	createTime := time.Now().UTC().UnixMilli()

	//TODO - Insert Manga
	create := pgdb.CreateMangaParams{
		MangaID: key,
		Title:   u.Title,
		TitleEn: pgtype.Text{
			String: u.TitleEnglish,
			Valid:  utils.StrIsEmpty(u.TitleEnglish)},
		Synonyms: u.Synonyms,
		Type:     u.Type,
		Country:  u.Country,
		Status: pgtype.Text{
			String: u.Status,
			Valid:  utils.StrIsEmpty(u.Status)},
		CreatedAt: createTime,
		UpdatedAt: createTime}

	mg, err := qtx.CreateManga(c, create)
	if err != nil {
		if err == pgsql.ErrRecordNotFound {
			return dtos.CreatedMangaResponse{MangaID: 0, Title: ""}, nil
		}
		tx.Rollback(c)
		rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return dtos.CreatedMangaResponse{}, err
	}
	//TODO - Created Genre when non exist
	for _, str := range u.Genres {
		key, err := flake.GenerateID(1, 3, 1)
		if err != nil {
			rp.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
			return dtos.CreatedMangaResponse{}, err
		}

		cgp := pgdb.CreateGenreParams{
			GenreID:   key,
			Title:     str,
			CreatedAt: createTime,
			UpdatedAt: createTime,
		}
		err = qtx.CreateGenre(c, cgp)
		if err != nil {
			tx.Rollback(c)
			rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
			return dtos.CreatedMangaResponse{}, err
		}
		time.Sleep(1 * time.Second)
	}

	//TODO - Find Genre ID by Tittle
	genreIds, err := qtx.FindGenreByTitle(c, u.Genres)
	if err != nil {
		tx.Rollback(c)
		rp.Log.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return dtos.CreatedMangaResponse{}, err
	}

	//TODO - Insert Manga Genres
	for _, ids := range genreIds {
		key, err := flake.GenerateID(1, 4, 1)
		if err != nil {
			tx.Rollback(c)
			rp.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
			return dtos.CreatedMangaResponse{}, err
		}
		cmgp := pgdb.CreateMangaGenreParams{
			MgID:      key,
			MangaID:   mg.MangaID,
			GenreID:   ids,
			UpdatedAt: createTime,
			CreatedAt: createTime,
		}
		err = qtx.CreateMangaGenre(c, cmgp)
		if err != nil {
			tx.Rollback(c)
			rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
			return dtos.CreatedMangaResponse{}, err
		}
		time.Sleep(1 * time.Second)
	}

	//TODO - Insert Manga Detail
	arg := pgdb.CreateMangaDetailParams{
		DetailID: mg.MangaID,
		Published: pgtype.Text{
			String: u.Published,
			Valid:  utils.StrIsEmpty(u.Published)},
		Authors: u.Authors,
		Artist:  u.Artists,
		Summary: pgtype.Text{
			String: u.Summary,
			Valid:  utils.StrIsEmpty(u.Summary)},
		CreatedAt: createTime,
		UpdatedAt: createTime,
	}
	err = qtx.CreateMangaDetail(c, arg)
	if err != nil {
		tx.Rollback(c)
		rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return dtos.CreatedMangaResponse{}, err
	}

	//TODO - Insert Manga Score
	arg2 := pgdb.CreateMangaScoreParams{
		ScoreID:   mg.MangaID,
		Score:     utils.FloatToPgNum(u.Score),
		UpdatedAt: createTime,
		CreatedAt: createTime,
	}
	err = qtx.CreateMangaScore(c, arg2)
	if err != nil {
		tx.Rollback(c)
		rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return dtos.CreatedMangaResponse{}, err
	}
	tx.Commit(c)
	return dtos.CreatedMangaResponse{MangaID: mg.MangaID, Title: mg.Title}, nil
}

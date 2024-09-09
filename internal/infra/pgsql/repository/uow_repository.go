package pgsql

import (
	"context"
	"fmt"
	"manga/config"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/infra/pgsql"
	"manga/internal/infra/pgsql/pgdb"
	"manga/internal/infra/rabbitmq"
	"manga/pkg"
	"manga/pkg/logging"
	"manga/pkg/utils"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meilisearch/meilisearch-go"
	"github.com/minio/minio-go/v7"
)

type uowRepository struct {
	q      *pgdb.Queries
	pgx    *pgxpool.Pool
	Log    logging.Logger
	cfg    *config.Config
	minioC *minio.Client
	rmq    *rabbitmq.Task
	ml     meilisearch.ServiceManager
}

func NewUOWRepository(pg *pkg.Postgres, cfg *config.Config, minioC *minio.Client, rmq *rabbitmq.Task, ml meilisearch.ServiceManager, Log logging.Logger) domain.IUOWRepository {
	return &uowRepository{
		pgx:    pg.Pool,
		Log:    Log,
		q:      pgdb.New(pg.Pool),
		cfg:    cfg,
		minioC: minioC,
		rmq:    rmq,
		ml:     ml,
	}
}

func (rp *uowRepository) CreateManga(c context.Context, u dtos.CreateManga) (dtos.CreateManga, error) {
	tx, err := rp.pgx.Begin(c)
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.CreateTx, err.Error(), nil)
		return dtos.CreateManga{}, err
	}

	qtx := rp.q.WithTx(tx)

	// - Generated Snow Flake ID
	key, err := pkg.GenerateID(1, 2, 1)
	if err != nil {
		rp.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
		tx.Rollback(c)
		return dtos.CreateManga{}, err
	}
	// - Generate Time
	createTime := time.Now().UTC().UnixMilli()

	// - Insert Manga
	create := pgdb.CreateMangaParams{
		MangaID: key,
		Title:   u.Title,
		Titles:  u.Titles,
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
			return dtos.CreateManga{}, nil
		}
		tx.Rollback(c)
		rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return dtos.CreateManga{}, err
	}

	// - Created Genre when non exist
	for _, str := range u.Genres {
		key, err := pkg.GenerateID(1, 3, 1)
		if err != nil {
			rp.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
			return dtos.CreateManga{}, err
		}

		cgp := pgdb.CreateGenreParams{
			GenreID:   key,
			Title:     str.Name,
			CreatedAt: createTime,
			UpdatedAt: createTime,
		}
		err = qtx.CreateGenre(c, cgp)
		if err != nil {
			tx.Rollback(c)
			rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
			return dtos.CreateManga{}, err
		}
		time.Sleep(1 * time.Millisecond)
	}

	// Find Genre ID by Tittle
	genreIds, err := qtx.FindGenreByTitle(c, dtos.MguMap(u.Genres,"Name"))
	if err != nil {
		tx.Rollback(c)
		rp.Log.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return dtos.CreateManga{}, err
	}

	// Insert Manga Genres
	for _, ids := range genreIds {
		key, err := pkg.GenerateID(1, 4, 1)
		if err != nil {
			tx.Rollback(c)
			rp.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
			return dtos.CreateManga{}, err
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
			return dtos.CreateManga{}, err
		}
		time.Sleep(1 * time.Millisecond)
	}

	// Insert Manga Detail
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
		return dtos.CreateManga{}, err
	}

	// Insert Manga Score
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
		return dtos.CreateManga{}, err
	}
	tx.Commit(c)
	u.MangaID = utils.Int64ToStr(mg.MangaID)
	return u, nil
}

func (rp *uowRepository) UpdateCover(c context.Context, u dtos.CreateManga) error{
	// Run the background goroutine to upload the image
	go UploadProcess(rp, u)

	return nil
}

func UploadProcess(rp *uowRepository, cm dtos.CreateManga)(bool,error ){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	success := true
	// Upload Cover Full
	
	filePath := fmt.Sprintf("%s/%s.webp", cm.MangaID, cm.MangaID)
	success, err := UploadToStorage(ctx, rp.cfg, rp.minioC, cm.Cover.CoverDetail, rp.cfg.Minio.Bucket1, filePath)
	if err != nil {
		fmt.Print(err.Error())

	}
	// Upload Cover Thumb
	thumbPath := fmt.Sprintf("%s/%s_thumb.webp", cm.MangaID, cm.MangaID)
	success, err = UploadToStorage(ctx, rp.cfg, rp.minioC, cm.Cover.Thumbnail, rp.cfg.Minio.Bucket1, thumbPath)
	if err != nil {
		fmt.Print(err.Error())

	}

	createTime := time.Now().UTC().UnixMilli()
	arg := pgdb.CreateMangaCoverParams{
		CoverID:     utils.StrToInt64(cm.MangaID),
		CoverDetail: pgtype.Text{String: filePath, Valid: utils.StrIsEmpty(filePath)},
		Thumbnail:   pgtype.Text{String: thumbPath, Valid: utils.StrIsEmpty(thumbPath)},
		UpdatedAt:   createTime,
		CreatedAt:   createTime,
	}
	_, err = rp.q.CreateMangaCover(ctx, arg)
	if err != nil {
		rp.Log.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		success = false

	}

	if !success {
		rp.Log.Error(logging.Rabbit,logging.Publish,fmt.Sprintf("Failed to publish to RabbitMQ: %v", err),nil)
		return false, pkg.NewErrorf(pkg.ErrorCodeUnknown,fmt.Sprintf("Failed to publish to RabbitMQ: %v", err))
	
	} 
	idxMg:= dtos.IndexedManga{
		ID: cm.MangaID,
		Title: cm.Title,
		Titles: cm.Titles,
		Synonyms: cm.Synonyms,
		Type: cm.Type,
		Genres: cm.Genres,
	   Cover: thumbPath,
	   Status: cm.Status,
	   Score: cm.Score,
	}
	  err = rp.rmq.CreatedManga(ctx,idxMg)
	  if err != nil{
		rp.Log.Error(logging.Rabbit,logging.Publish,fmt.Sprintf("Failed to publish to MeiliSearch:: %v", err),nil)
	  }
	  return true, nil
}

func UploadToStorage(ctx context.Context, cfg *config.Config, minioC *minio.Client, source string, bucket string, filePath string) (bool, error) {

	if !strings.Contains(source, ".webp") {

		url := fmt.Sprintf("%s:%s/unsafe/filters:format(webp)/%s", cfg.Imagor.Host, cfg.Imagor.Port, source)
		source = url
	}

	res, err := utils.GetFileRequest(source)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	_, err = minioC.PutObject(ctx, bucket, filePath, res.Body, res.ContentLength, minio.PutObjectOptions{
		ContentType: res.Header.Get("Content-Type"), // You can specify the content type
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

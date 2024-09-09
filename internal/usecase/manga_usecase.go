package usecase

import (
	"context"
	"fmt"

	"manga/config"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"manga/internal/domain/params"
	"manga/internal/infra/meili"
	"manga/pkg/logging"
	"manga/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type IMangaUsecase interface {
	Create(c context.Context, m dtos.CreateManga) (dtos.CreatedMangaResponse, error)
	FindById(c context.Context, id string) (dtos.Manga, error)
	Find(c context.Context,search params.SearchParams)(any,any,error)
}
type mangaUsecase struct {
	conf           *config.Config
	uowRepo        domain.IUOWRepository
	mangaRepo      domain.IMangaRepository
	meiliRepo      *meili.Manga
	contextTimeout time.Duration
	redis          *redis.Client
	log            logging.Logger
}

func NewMangaUsecase(uowRepository domain.IUOWRepository, mangaRepo domain.IMangaRepository, timeout time.Duration, redis *redis.Client, cfg *config.Config,meiliRepo *meili.Manga, log logging.Logger) IMangaUsecase {
	return &mangaUsecase{
		uowRepo:        uowRepository,
		mangaRepo:      mangaRepo,
		contextTimeout: timeout,
		redis:          redis,
		conf:           cfg,
		log:            log,
		meiliRepo: meiliRepo,
	}
}

func (mu *mangaUsecase) Create(c context.Context, m dtos.CreateManga) (dtos.CreatedMangaResponse, error) {
	ctx, cancel := context.WithTimeout(c, 20*time.Second)
	defer cancel()

	if m.Method == "indirect" {
		res, err := mu.uowRepo.CreateManga(ctx, m)
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		return dtos.CreatedMangaResponse{MangaID: res.MangaID,Title: res.Title}, nil
	}
	switch m.Source {
	// Generated Post by My Anime List
	case "myanimelist":
		url := fmt.Sprintf("%s/manga/%s", mu.conf.Source.MyAnimeList, m.MangaID)
		//Http requested for Jikan Api
		mal, err := utils.GetHttpRequest[models.JikanMangaByID](url, mu.log)
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		//Created manga
		mg, err := mu.uowRepo.CreateManga(ctx, dtos.CreateByMyAnimeList(mal))
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		var mangaIdx dtos.CreateManga
        genre, err := mu.mangaRepo.FindMangaGenre(c, mg.MangaID)
		if err != nil {
		return dtos.CreatedMangaResponse{}, err
		}
		mangaIdx = mg
		mangaIdx.Genres = dtos.ToGenreFromPGRB(genre)
		err  = mu.uowRepo.UpdateCover(ctx,mangaIdx)
        if err != nil{
			return dtos.CreatedMangaResponse{},err
		}
		return dtos.CreatedMangaResponse{MangaID: mg.MangaID,Title: mg.Title}, nil
		// Generated Post by My Anime List
	case "mangaupdate":
		url := fmt.Sprintf("%s/v1/series/%s", mu.conf.Source.MangaUpdate, m.MangaID)
		mgu, err := utils.GetHttpRequest[models.MangaUpdateByID](url, mu.log)
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		mg, err := mu.uowRepo.CreateManga(ctx, dtos.CreateByMangaUpdate(mgu))
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		var mangaIdx dtos.CreateManga
        genre, err := mu.mangaRepo.FindMangaGenre(c, mg.MangaID)
		if err != nil {
		return dtos.CreatedMangaResponse{}, err
		}
		mangaIdx = mg
		mangaIdx.Genres = dtos.ToGenreFromPGRB(genre)
		err  = mu.uowRepo.UpdateCover(ctx,mangaIdx)
        if err != nil{
			return dtos.CreatedMangaResponse{},err
		}
		return dtos.CreatedMangaResponse{MangaID: mg.MangaID,Title: mg.Title}, nil
	default:
		return dtos.CreatedMangaResponse{}, utils.ErrorWrapper("source not found")
	}
}

func (mu *mangaUsecase) FindById(c context.Context, id string) (dtos.Manga, error) {
	ctx, cancel := context.WithTimeout(c, mu.contextTimeout)
	defer cancel()
	mg, err := mu.mangaRepo.FindbyId(ctx, id)
	if err != nil {
		return dtos.Manga{}, err
	}
	genre, err := mu.mangaRepo.FindMangaGenre(c, id)
	if err != nil {
		return dtos.Manga{}, err
	}
	return dtos.ToManga(mg, genre), nil
}
func (mu *mangaUsecase) Find(c context.Context,search params.SearchParams)(any,any,error){
    ctx, cancel := context.WithTimeout(c, mu.contextTimeout)
	defer cancel()

	res,page,err:= mu.meiliRepo.Search(ctx,search)
	if err != nil{
		return nil,nil,err
	}

	return res,page,nil
}

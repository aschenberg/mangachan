package usecase

import (
	"context"
	"fmt"
	"reflect"

	"manga/config"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"manga/pkg/logging"
	"manga/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type IMangaUsecase interface {
	Create(c context.Context, m dtos.CreateManga) (dtos.CreatedMangaResponse, error)
	FindById(c context.Context, id string) (dtos.Manga, error)
}
type mangaUsecase struct {
	conf           *config.Config
	uowRepo        domain.IUOWRepository
	mangaRepo      domain.IMangaRepository
	contextTimeout time.Duration
	redis          *redis.Client
	log            logging.Logger
}

func NewMangaUsecase(uowRepository domain.IUOWRepository, mangaRepo domain.IMangaRepository, timeout time.Duration, redis *redis.Client, cfg *config.Config, log logging.Logger) IMangaUsecase {
	return &mangaUsecase{
		uowRepo:        uowRepository,
		mangaRepo:      mangaRepo,
		contextTimeout: timeout,
		redis:          redis,
		conf:           cfg,
		log:            log,
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
		return res, nil
	}
	switch m.Source {
	//TODO - Generated Post by My Anime List
	case "myanimelist":
		url := fmt.Sprintf("%s/manga/%s", mu.conf.Source.MyAnimeList, m.MangaID)
		mal, err := utils.GetHttpRequest[models.JikanMangaByID](url, mu.log)
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		mg := dtos.CreateManga{
			Title:        mal.Data.Title,
			TitleEnglish: mal.Data.TitleEnglish,
			Synonyms:     mal.Data.TitleSynonyms,
			Type:         mal.Data.Type,
			Published:    mal.Data.Published.String,
			Country:      switchCountry(mal.Data.Type),
			Status:       mal.Data.Status,
			Authors:      addName(mal.Data.Authors),
			Artists:      addName(mal.Data.Authors),
			Genres:       addName(mal.Data.Genres),
			Themes:       addName(mal.Data.Themes),
			Demographic:  addName(mal.Data.Demographics),
			Summary:      mal.Data.Synopsis,
			Score:        mal.Data.Score,
			Cover: models.MangaCover{
				CoverDetail: mal.Data.Images.Webp.LargeImageURL,
				Thumbnail:   mal.Data.Images.Webp.ImageURL,
			},
		}
		res, err := mu.uowRepo.CreateManga(ctx, mg)
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		return res, nil
	case "mangaupdate":

		url := fmt.Sprintf("%s/v1/series/%s", mu.conf.Source.MangaUpdate, m.MangaID)
		mgu, err := utils.GetHttpRequest[models.MangaUpdateByID](url, mu.log)
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		mg := dtos.CreateManga{
			Title:        mgu.Title,
			TitleEnglish: mgu.Title,
			Synonyms:     mguMap(mgu.Associated, "Title"),
			Type:         mgu.Type,
			Published:    mgu.Year,
			Country:      switchCountry(mgu.Type),
			Status:       mgu.Status,
			Authors:      mguMap(mgu.Authors, "Name"),
			Artists:      mguMap(mgu.Authors, "Name"),
			Genres:       mguMap(mgu.Genres, "Genre"),
			Themes:       mguMap(mgu.Categories, "Category"),
			Demographic:  []string{},
			Summary:      mgu.Description,
			Score:        float64(mgu.BayesianRating),
			Cover: models.MangaCover{
				CoverDetail: mgu.Image.URL.Original,
				Thumbnail:   mgu.Image.URL.Original,
			},
		}
		res, err := mu.uowRepo.CreateManga(ctx, mg)
		if err != nil {
			return dtos.CreatedMangaResponse{}, err
		}
		return res, nil
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

func addName(vals []models.MALSubJson) []string {
	var valStr []string
	for _, a := range vals {
		valStr = append(valStr, a.Name)
	}
	return valStr
}

func mguMap(vals interface{}, field string) []string {
	valsValue := reflect.ValueOf(vals)
	if valsValue.Kind() != reflect.Slice {
		panic("mguMap: input must be a slice")
	}

	var valStr []string
	for i := 0; i < valsValue.Len(); i++ {
		elem := valsValue.Index(i)
		field := elem.FieldByName(field)
		if !field.IsValid() {
			panic("mguMap: struct does not have a Title field")
		}
		valStr = append(valStr, field.String())
	}
	return valStr
}

func switchCountry(val string) string {
	switch val {
	case "Manga":
		return "jp"
	case "manga":
		return "jp"
	case "Manhwa":
		return "kr"
	case "manhwa":
		return "kr"
	case "Manhua":
		return "cn"
	case "manhua":
		return "cn"
	default:
		return "unknown"
	}
}

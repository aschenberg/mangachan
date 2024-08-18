package usecase

import (
	"context"

	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"time"

	"github.com/redis/go-redis/v9"
)

type IMangaUsecase interface {
	Create(c context.Context, m dtos.CreateManga) (dtos.CreatedMangaResponse, error)
}
type mangaUsecase struct {
	uowRepository  domain.IUOWRepository
	contextTimeout time.Duration
	redis          *redis.Client
}

func NewMangaUsecase(uowRepository domain.IUOWRepository, timeout time.Duration, redis *redis.Client) IMangaUsecase {
	return &mangaUsecase{
		uowRepository:  uowRepository,
		contextTimeout: timeout,
		redis:          redis,
	}
}

func (mu *mangaUsecase) Create(c context.Context, m dtos.CreateManga) (dtos.CreatedMangaResponse, error) {
	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()
	mg := dtos.CreateManga{
		Title:        m.Title,
		TitleEnglish: m.TitleEnglish,
		Synonyms:     m.Synonyms,
		Type:         m.Type,
		Published:    m.Published,
		Country:      m.Country,
		Status:       m.Status,
		Authors:      m.Authors,
		Artists:      m.Artists,
		Genres:       m.Genres,
		Themes:       m.Themes,
		Demographic:  m.Demographic,
		Summary:      m.Summary,
		Score:        m.Score,
		Cover:        m.Cover,
	}
	return mu.uowRepository.CreateManga(ctx, mg)
}

// func (mu *mangaUsecase) FindById(c context.Context, id string) (models.Manga, error) {
// 	ctx, cancel := context.WithTimeout(c, mu.contextTimeout)
// 	defer cancel()
// 	err := mu.redis.Set(ctx, id, id, 10*time.Second).Err()
// 	if err != nil {
// 		return models.Manga{}, err
// 	}
// 	return mu.mangaRepository.FindById(ctx, id)
// }

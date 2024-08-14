package usecase

import (
	"context"

	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
	"manga/pkg"
	"manga/pkg/fileutils"
	"manga/pkg/flake"
	"manga/pkg/imagekio"
	"manga/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type mangaUsecase struct {
	mangaRepository domain.MangaRepository
	contextTimeout  time.Duration
	redis           *redis.Client
}

func NewMangaUsecase(mangaRepository domain.MangaRepository, timeout time.Duration, redis *redis.Client) domain.MangaUsecase {
	return &mangaUsecase{
		mangaRepository: mangaRepository,
		contextTimeout:  timeout,
		redis:           redis,
	}
}

func (mu *mangaUsecase) Create(c context.Context, manga dtos.CreateManga) error {
	var err error
	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()
	newId, err := flake.GenerateID(1, 2, 1)
	nowTime := time.Now().UTC()
	imgID, err := flake.GenerateID(2, 1, 1)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "error generate spaceflake")
	}

	result, err := imagekio.UploadSingle(ctx, manga.Cover, newId.Hex(), utils.Int64ToStr(imgID), "jpg")
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "error send image")
	}

	err = fileutils.SaveFileFromURL(result.Data.Url, newId.Hex())
	if err != nil {
		return err
	}

	newManga := models.Manga{
		Manga_ID:     newId,
		Title:        manga.Title,
		TitleEnglish: manga.TitleEnglish,
		Synonyms:     manga.Synonyms,
		Type:         manga.Type,

		Country:   manga.Country,
		Published: manga.Published,
		Status:    manga.Status,
		Authors:   manga.Authors,
		Artists:   manga.Artists,

		Score: manga.Score,

		Summary:   manga.Summary,
		UpdatedAt: nowTime,
		CreatedAt: nowTime,
	}
	return mu.mangaRepository.Create(ctx, newManga)
}

func (mu *mangaUsecase) FindById(c context.Context, id string) (models.Manga, error) {
	ctx, cancel := context.WithTimeout(c, mu.contextTimeout)
	defer cancel()
	print(id)
	err := mu.redis.Set(ctx, id, id, 10*time.Second).Err()
	if err != nil {
		return models.Manga{}, err
	}
	return mu.mangaRepository.FindById(ctx, id)
}

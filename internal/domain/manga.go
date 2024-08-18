package domain

import (
	"context"
	"manga/internal/domain/dtos"
	"manga/internal/domain/models"
)

type IMangaRepository interface {
	Create(c context.Context, manga models.Manga) error
	FindById(c context.Context, id string) (models.Manga, error)
}

type MangaUsecase interface {
	Create(c context.Context, manga dtos.CreateManga) error
	// FindById(c context.Context, id string) (models.Manga, error)
}

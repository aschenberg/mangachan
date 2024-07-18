package domain

import (
	"context"
	"manga/domain/dtos"
	"manga/domain/models"
)

type MangaRepository interface {
	Create(c context.Context, manga models.Manga) error
	FindById(c context.Context, id string) (models.Manga, error)
}

type MangaUsecase interface {
	Create(c context.Context, manga dtos.CreateManga) error
	FindById(c context.Context, id string) (models.Manga, error)
}

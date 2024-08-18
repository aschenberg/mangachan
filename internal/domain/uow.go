package domain

import (
	"context"
	"manga/internal/domain/dtos"
)

type IUOWRepository interface {
	CreateManga(c context.Context, u dtos.CreateManga) (dtos.CreatedMangaResponse, error)
	// FindById(c context.Context, id string) (models.Manga, error)
}

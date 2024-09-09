package domain

import (
	"context"
	"manga/internal/domain/dtos"
)

type IUOWRepository interface {
	CreateManga(c context.Context, u dtos.CreateManga) (dtos.CreateManga, error)
	UpdateCover(c context.Context, u dtos.CreateManga) error
}

package domain

import (
	"context"
	"manga/internal/domain/models"
)

type ChapterRepository interface {
	Create(c context.Context, chapter models.Chapter) error
}

type ChapterUsecase interface {
}

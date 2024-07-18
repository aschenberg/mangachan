package domain

import (
	"context"
	"manga/domain/models"
)

type ChapterRepository interface {
	Create(c context.Context, chapter models.Chapter) error
}

type ChapterUsecase interface {
}

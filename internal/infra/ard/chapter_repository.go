package ard

import (
	"context"
	"manga/db"
	"manga/internal/domain"
	"manga/internal/domain/models"
)

type chapterRepository struct {
	database   db.MongoDB
	collection string
}

func NewChapterRepository(db db.MongoDB, collection string) domain.ChapterRepository {
	return &chapterRepository{
		database:   db,
		collection: collection,
	}
}

func (cr *chapterRepository) Create(c context.Context, chapter models.Chapter) error {
	collection := cr.database.Collection(cr.collection)

	_, err := collection.InsertOne(c, chapter)
	if err != nil {
		return err
	}

	return nil
}

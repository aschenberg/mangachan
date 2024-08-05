package mongodb

import (
	"context"
	"manga/db"
	"manga/internal/domain"
	"manga/internal/domain/models"

	"manga/pkg"

	"go.mongodb.org/mongo-driver/bson"
)

type mangaRepository struct {
	database   db.MongoDB
	collection string
}

func NewMangaRepository(db db.MongoDB, collection string) domain.MangaRepository {
	return &mangaRepository{
		database:   db,
		collection: collection,
	}
}

func (cr *mangaRepository) Create(c context.Context, manga models.Manga) error {
	collection := cr.database.Collection(cr.collection)

	_, err := collection.InsertOne(c, manga)
	if err != nil {
		return err
	}

	return nil
}

func (cr *mangaRepository) FindById(c context.Context, id string) (models.Manga, error) {
	collection := cr.database.Collection(models.CollectionManga)
	var result models.Manga

	err := collection.FindOne(c, bson.M{"_id": pkg.StrToObjectID(id)}).Decode(&result)
	if err != nil {
		return models.Manga{}, err
	}

	return result, nil
}

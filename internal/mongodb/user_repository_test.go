package mongodb

import (
	"context"
	"errors"
	"testing"

	"manga/db/mocks"
	"manga/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreate(t *testing.T) {

	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}

	collectionName := models.CollectionUser

	mockUser := models.User{
		ID:    primitive.NewObjectID(),
		Name:  "Test",
		Email: "test@gmail.com",
	}

	mockEmptyUser := models.User{}
	mockUserID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {

		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.User")).Return(mockUserID, nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := NewUserRepository(databaseHelper, collectionName)

		_, err := ur.Create(context.Background(), mockUser)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.User")).Return(mockEmptyUser, errors.New("Unexpected")).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := NewUserRepository(databaseHelper, collectionName)

		_, err := ur.Create(context.Background(), mockEmptyUser)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})

}

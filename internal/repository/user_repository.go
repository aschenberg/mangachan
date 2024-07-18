package repository

import (
	"context"
	"fmt"

	"manga/db"
	"manga/domain"

	"manga/domain/models"
	"manga/pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	database   db.MongoDB
	collection string
}

func NewUserRepository(db db.MongoDB, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (ur *userRepository) Create(c context.Context, user models.User) (models.User, error) {
	collection := ur.database.Collection(ur.collection)

	result, err := collection.InsertOne(c, user)
	if err != nil {
		return models.User{}, err
	}

	var createdUser models.User
	err = collection.FindOne(c, bson.M{"_id": result.InsertedID}).Decode(&createdUser)
	if err != nil {
		return models.User{}, pkg.NewErrorf(pkg.ErrorCodeUnknown, "failed to retrieve created user: %v", err)
	}
	return createdUser, nil
}

func (ur *userRepository) Fetch(c context.Context) ([]models.User, error) {
	collection := ur.database.Collection(ur.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []models.User

	err = cursor.All(c, &users)
	if users == nil {
		return []models.User{}, err
	}

	return users, err
}

func (ur *userRepository) GetByAppId(c context.Context, appId string) (models.User, error) {
	collection := ur.database.Collection(ur.collection)
	var user models.User
	err := collection.FindOne(c, bson.M{"app_id": appId}).Decode(&user)
	return user, err
}

func (ur *userRepository) GetByID(c context.Context, id string) (models.User, error) {
	collection := ur.database.Collection(ur.collection)

	var user models.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return user, err
}

func (ur *userRepository) FindByRefreshToken(c context.Context, token string) (models.User, error) {
	collection := ur.database.Collection(ur.collection)
	filter := bson.M{"refresh_token": token}
	var result models.User
	err := collection.FindOne(c, filter).Decode(&result)
	if err != nil {
		return models.User{}, err
	}
	return result, nil
}
func (ur *userRepository) UpdateRefreshToken(c context.Context, id string, token string) error {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"_id": pkg.StrToObjectID(id)}
	update := bson.M{"$set": bson.M{"refresh_token": token}}

	updateResult, err := collection.UpdateOne(c, filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return err
	}
	return nil
}

// func (ur *userRepository) FetchByName(c context.Context, name string, sortField string, sortOrder int, page int64, pageSize int64, rsid string) ([]models.User, int64, error) {
// 	collection := ur.database.Collection(ur.collection)

// 	var filter bson.D
// 	if rsid == "" {
// 		if name == "" {
// 			filter = bson.D{} // Empty filter matches all documents
// 		} else {
// 			filter = bson.D{
// 				{Key: "name", Value: bson.D{
// 					{Key: "$regex", Value: primitive.Regex{Pattern: name, Options: "i"}},
// 				}},
// 			} // Build the filter to find documents by name
// 		}
// 	} else {
// 		if name == "" {
// 			filter = bson.D{
// 				{Key: "place", Value: utils.StrToObjectID(rsid)},
// 			} // Rsid filter matches all documents
// 		} else {
// 			filter = bson.D{
// 				{Key: "name", Value: bson.D{
// 					{Key: "place", Value: utils.StrToObjectID(rsid)},
// 					{Key: "$regex", Value: primitive.Regex{Pattern: name, Options: "i"}},
// 				}},
// 			} // Build the filter to find documents by name
// 		}
// 	}

// 	// Count total documents
// 	totalDocs, err := collection.CountDocuments(c, filter)
// 	if err != nil {
// 		return []domain.User{}, 0, err
// 	}

// 	// Calculate total pages
// 	totalPages := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))

// 	// Options for sorting and pagination
// 	findOptions := options.Find()
// 	findOptions.SetSort(bson.D{{Key: sortField, Value: -1}})
// 	findOptions.SetSkip((page - 1) * pageSize)
// 	findOptions.SetLimit(pageSize)
// 	cursor, err := collection.Find(c, filter, findOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(c)

// 	var results []domain.User
// 	if err = cursor.All(c, &results); err != nil {
// 		utils.ErrorWrapper("error", err)
// 	}

// 	return results, totalPages, err
// }

func (ur *userRepository) DeleteByID(c context.Context, id string) error {
	collection := ur.database.Collection(ur.collection)

	result, err := collection.DeleteOne(c, bson.M{"_id": pkg.StrToObjectID(id)})
	if err != nil {
		return err
	}
	if result == 0 {
		fmt.Println("No document was deleted")
	} else {
		fmt.Printf("Deleted %d document(s)\n", result)
	}
	return nil
}

package ard

// import (
// 	"context"
// 	"fmt"

// 	"manga/internal/domain"
// 	"manga/internal/domain/dtos"
// 	"manga/internal/domain/models"

// 	"github.com/arangodb/go-driver/v2/arangodb"
// 	"github.com/arangodb/go-driver/v2/arangodb/shared"
// )

// type userRepository struct {
// 	db  arangodb.Database
// 	col string
// }

// func NewUserRepository(db arangodb.Database, col string) domain.UserRepository {
// 	return &userRepository{
// 		db:  db,
// 		col: col,
// 	}
// }

// func (ur *userRepository) Create(c context.Context, user models.User) (models.User, error) {
// 	col, err := ur.db.GetCollection(c, ur.col, nil)
// 	if err != nil {
// 		return models.User{}, err
// 	}

// 	_, err = col.CreateDocument(c, user)
// 	if err != nil {
// 		return models.User{}, err
// 	}
// 	// user.ID = string(result.ID)
// 	// user.Key = result.Key

// 	return models.User{}, nil
// }

// func (ur *userRepository) Fetch(c context.Context) ([]models.User, error) {
// 	collection := ur.database.Collection(ur.collection)

// 	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
// 	cursor, err := collection.Find(c, bson.D{}, opts)

// 	if err != nil {
// 		return nil, err
// 	}

// 	var users []models.User

// 	err = cursor.All(c, &users)
// 	if users == nil {
// 		return []models.User{}, err
// 	}

// 	return users, err
// }

// func (ur *userRepository) UpdateBySubID(c context.Context, claim models.GoogleClaims) (models.User, error) {
// 	query := `FOR user IN @collection
// 			  FILTER user.app_id == @app_id
// 			  UPDATE user WITH { picture: @picture, given_name: @given_name, family_name: @family_name, name: @name } IN @collection`
// 	bindVars := map[string]interface{}{
// 		"collection":  ur.col,
// 		"app_id":      claim.Sub,
// 		"picture":     claim.Picture,
// 		"given_name":  claim.GivenName,
// 		"family_name": claim.FamilyName,
// 		"name":        claim.Name,
// 	}
// 	cursor, err := ur.db.Query(c, query, &arangodb.QueryOptions{BindVars: bindVars})
// 	defer cursor.Close()

// 	var updatedUser models.User

// 	_, err = cursor.ReadDocument(context.Background(), &updatedUser)
// 	if err != nil {
// 		return models.User{}, fmt.Errorf("failed to read updated document: %w", err)
// 	}
// 	return updatedUser, nil
// }

// func (ur *userRepository) IsExistBySubID(c context.Context, subId string) (bool, error) {
// 	query := `FOR user IN @@collection FILTER user.app_id == @subId LIMIT 1 RETURN user`
// 	cursor, err := ur.db.Query(c, query,
// 		&arangodb.QueryOptions{
// 			BindVars: map[string]interface{}{"@collection": ur.col, "subId": subId}})
// 	if err != nil {

// 		return false, fmt.Errorf("failed to execute query: %w", err)
// 	}
// 	defer cursor.Close()

// 	// Read the first document
// 	var res models.User
// 	_, err = cursor.ReadDocument(c, &res)
// 	if shared.IsNoMoreDocuments(err) {
// 		// No document found
// 		return false, nil
// 	} else if err != nil {
// 		return false, fmt.Errorf("failed to read document: %w", err)
// 	}
// 	return true, nil
// }

// func (ur *userRepository) CreateOrUpdate(c context.Context, claim models.GoogleClaims) (models.User, string, error) {

// 	query := `
// 	UPSERT { app_id: @app_id }
// 	INSERT {
// 		app_id: @app_id,
// 		email: @Email,
// 		picture: @Picture,
// 		role: @Role,
// 		given_name: @GivenName,
// 		family_name: @FamilyName,
// 		name: @Name
// 	}
// 	UPDATE {
// 		given_name: @GivenName,
// 		family_name: @FamilyName,
// 		name: @Name,
// 		picture: @Picture

// 	}
// 	IN @@collection
// 	RETURN { doc: NEW, type: OLD ? 'update' : 'insert' }`

// 	bindVars := map[string]interface{}{
// 		"@collection": ur.col,
// 		"app_id":      claim.Sub,
// 		"Email":       claim.Email,
// 		"Picture":     claim.Picture,
// 		"Role":        []string{"user"},
// 		"GivenName":   claim.GivenName,
// 		"FamilyName":  claim.FamilyName,
// 		"Name":        claim.Name,
// 	}
// 	opts := &arangodb.QueryOptions{BindVars: bindVars}
// 	cursor, err := ur.db.Query(c, query, opts)
// 	if err != nil {
// 		return models.User{}, "", fmt.Errorf("failed to execute query: %w", err)
// 	}
// 	defer cursor.Close()

// 	var result dtos.AuthResult

// 	_, err = cursor.ReadDocument(c, &result)
// 	if shared.IsNoMoreDocuments(err) {
// 		return models.User{}, "", nil
// 	} else if err != nil {
// 		return models.User{}, "", fmt.Errorf("failed to read document: %w", err)
// 	}
// 	return result.Doc, result.Type, nil

// }

//	func (ur *userRepository) FindByRefreshToken(c context.Context, token string) (models.User, error) {
//		collection := ur.database.Collection(ur.collection)
//		filter := bson.M{"refresh_token": token}
//		var result models.User
//		err := collection.FindOne(c, filter).Decode(&result)
//		if err != nil {
//			return models.User{}, err
//		}
//		return result, nil
//	}
// func (ur *userRepository) UpdateRefreshToken(c context.Context, subId string, token string) error {
// 	query := `FOR user IN @collection
// 			  FILTER user.app_id == @app_id
// 			  UPDATE user WITH { refresh_token: @token } IN @collection`
// 	_, err := ur.db.Query(c, query,
// 		&arangodb.QueryOptions{
// 			BindVars: map[string]interface{}{"app_id": subId, "token": token}})
// 	if err != nil {
// 		return fmt.Errorf("failed to execute query: %w", err)
// 	}

// 	return nil
// }

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

// func (ur *userRepository) DeleteByID(c context.Context, id string) error {
// 	collection := ur.database.Collection(ur.collection)

// 	result, err := collection.DeleteOne(c, bson.M{"_id": pkg.StrToObjectID(id)})
// 	if err != nil {
// 		return err
// 	}
// 	if result == 0 {
// 		fmt.Println("No document was deleted")
// 	} else {
// 		fmt.Printf("Deleted %d document(s)\n", result)
// 	}
// 	return nil
// }

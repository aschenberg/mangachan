package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionUser = "users"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	AppID      string             `bson:"app_id"`
	Email      string             `bson:"email"`
	Picture    string             `bson:"picture"`
	Role       []string           `bson:"role"`
	GivenName  string             `bson:"given_name"`
	FamilyName string             `bson:"family_name"`
	Name       string             `bson:"name"`
}

package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	Picture    string   `json:"picture"`
	Role       []string `json:"role"`
	GivenName  string   `json:"given_name"`
	FamilyName string   `json:"family_name"`
	jwt.MapClaims
}

type JwtCustomRefreshClaims struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	Picture    string   `json:"picture"`
	Role       []string `json:"role"`
	GivenName  string   `json:"given_name"`
	FamilyName string   `json:"family_name"`
	jwt.MapClaims
}

package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Role  []string `json:"role"`
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

type GoogleClaims struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Athash        string `json:"at_hash"`
	Hd            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
}

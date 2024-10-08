package pkg

import (
	"fmt"
	"manga/internal/domain/models"

	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user models.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).UnixMilli()
	now := time.Now().UnixMilli()
	claims := &models.JwtCustomClaims{
		Role: user.Role,
		MapClaims: jwt.MapClaims{
			"sub":  user.UserID,
			"name": user.GivenName + user.FamilyName,
			"iat":  now,
			"exp":  exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(user models.User, secret string, expiry int) (refreshToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).UnixMilli()
	now := time.Now().UnixMilli()
	claimsRefresh := &models.JwtCustomRefreshClaims{
		Role: user.Role,
		MapClaims: jwt.MapClaims{
			"sub":  user.UserID,
			"name": user.GivenName + user.FamilyName,
			"iat":  now,
			"exp":  exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	return claims["sub"].(string), nil
}

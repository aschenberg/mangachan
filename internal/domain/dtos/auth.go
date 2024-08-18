package dtos

import "manga/internal/domain/models"

type AuthRequest struct {
	IDToken string `json:"id_token"`
}

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Role         []string `json:"role"`
	Picture      string   `json:"picture"`
}

type AuthResult struct {
	Doc  models.User `json:"doc"`
	Type string      `json:"type"`
}

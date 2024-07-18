package dtos

type AuthRequest struct {
	IDToken     string `json:"idToken"`
	AccessToken string `json:"accessToken"`
}

type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Role         []string `json:"role"`
	Picture      string   `json:"picture"`
}

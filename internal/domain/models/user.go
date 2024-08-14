package models

const (
	CollectionUser = "users"
)

type User struct {
	UserID       string   `json:"user_id,omitempty"`
	AppID        string   `json:"app_id"`
	Email        string   `json:"email"`
	Picture      string   `json:"picture"`
	Role         []string `json:"role"`
	GivenName    string   `json:"given_name"`
	FamilyName   string   `json:"family_name"`
	Name         string   `json:"name"`
	RefreshToken string   `json:"refresh_token"`
}
type UserCreated struct {
	AppID        string   `json:"app_id"`
	Email        string   `json:"email"`
	Picture      string   `json:"picture"`
	Role         []string `json:"role"`
	GivenName    string   `json:"given_name"`
	FamilyName   string   `json:"family_name"`
	Name         string   `json:"name"`
	RefreshToken string   `json:"refresh_token"`
}

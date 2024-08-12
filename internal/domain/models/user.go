package models

const (
	CollectionUser = "users"
)

type User struct {
	ID           string   `json:"_id,omitempty"`
	Key          string   `json:"_key,omitempty"`
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

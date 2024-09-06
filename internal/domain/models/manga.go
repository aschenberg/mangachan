package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionManga = "manga"
)

type Manga struct {
	Manga_ID     int64     `json:"manga_id"`
	Title        string    `json:"title"`
	TitleEnglish string    `json:"title_en"`
	Synonyms     []string  `json:"synonyms"`
	CoverID      int64     `json:"cover_id"`
	Type         string    `json:"type"`
	Country      string    `json:"country"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type MangaDetail struct {
	Detail_ID int64     `json:"detail_id"`
	Published string    `json:"published"`
	Authors   []string  `json:"authors"`
	Artists   []string  `json:"artists"`
	Summary   string    `json:"summary"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
type MangaScrore struct {
	ScoreID int64   `json:"score_id"`
	Score   float64 `json:"score"`
}

type MangaCover struct {
	CoverID     string   `json:"cover_id"`
	CoverDetail string   `json:"cover_detail"`
	Thumbnail   string   `json:"thumbnail"`
	Extra       []string `json:"extra,omitempty"`
	UpdatedAt   int64    `json:"updated_at,omitempty"`
	CreatedAt   int64    `json:"created_at,omitempty"`
}

type Genre struct {
	GenreID string `json:"genre_id"`
	Name    string `json:"name"`
}

type Theme struct {
	ThemeID int    `json:"theme_id"`
	Name    string `json:"name"`
}

type Demographic struct {
	DemographicID int    `json:"demographic_id"`
	Name          string `json:"name"`
}

type MangaComment struct {
	ID        primitive.ObjectID `json:"_id"`
	MangaID   primitive.ObjectID `json:"manga_id"`
	UserId    primitive.ObjectID `json:"user_id"`
	Comment   string             `json:"comment"`
	CreatedAt time.Time          `json:"created_at"`
}

type ImageCover struct {
	LocalUrl  string `json:"local_url"`
	RemoteUrl string `json:"remote_url"`
}

type Rating struct {
	ID      primitive.ObjectID `json:"_id,omitempty"`
	MangaID primitive.ObjectID `json:"manga_id"`
	UserID  primitive.ObjectID `json:"user_id"`
	Score   int                `json:"score"`
}

type Result struct {
	Status string `json:"status"`
}

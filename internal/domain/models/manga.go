package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionManga = "manga"
)

type Manga struct {
	ID           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	TitleEnglish string             `bson:"title_en"`
	Synonyms     []string           `bson:"synonyms"`
	Cover        ImageCover         `bson:"cover"`
	Type         string             `bson:"type"`
	Country      string             `bson:"country"`
	Published    string             `bson:"published"`
	Status       string             `bson:"status"`
	Authors      []string           `bson:"authors"`
	Artists      []string           `bson:"artists"`
	Genres       []string           `bson:"genres"`
	Score        float64            `bson:"score"`
	Themes       []string           `bson:"themes"`
	Demographic  []string           `bson:"demographic"`
	Summary      string             `bson:"summary"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	CreatedAt    time.Time          `bson:"created_at"`
}

type MangaComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	MangaID   primitive.ObjectID `bson:"manga_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Comment   string             `bson:"comment"`
	CreatedAt time.Time          `bson:"created_at"`
}

type ImageCover struct {
	LocalUrl  string `bson:"local_url"`
	RemoteUrl string `bson:"remote_url"`
}

type Rating struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	MangaID primitive.ObjectID `bson:"manga_id"`
	UserID  primitive.ObjectID `bson:"user_id"`
	Score   int                `bson:"score"`
}

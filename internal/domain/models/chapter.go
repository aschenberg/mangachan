package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chapter struct {
	ID         primitive.ObjectID `bson:"_id"`
	MangaID    primitive.ObjectID `bson:"manga_id"`
	MangaTitle string             `bson:"manga_title"`
	Number     string             `bson:"number"`
	Image      primitive.ObjectID `bson:"image_id"`
	CreatedAt  time.Time          `bson:"created_at"`
}

type ChapterLanguage struct {
	ID         primitive.ObjectID `bson:"_id"`
	MangaID    primitive.ObjectID `bson:"manga_id"`
	MangaTitle string             `bson:"manga_title"`
	Title      string             `bson:"title"`
	Source     string             `bson:"source"`
	SourceUrl  string             `bson:"source_url"`
	Lang       string             `bson:"lang"`
	Uploader   Uploader           `bson:"uploader"`
	Viewer     primitive.ObjectID `bson:"viewer_id"`
	Translator Translator         `bson:"translator"`
	CreatedAt  time.Time          `bson:"created_at"`
}

type Translator struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Uploader struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type ChapterImage struct {
	ID        primitive.ObjectID `bson:"_id"`
	MangaID   primitive.ObjectID `bson:"manga_id"`
	ChapterID primitive.ObjectID `bson:"chapter_id"`
	Images    []Image            `bson:"images"`
}
type Image struct {
	Index     int       `bson:"idx"`
	Url       string    `bson:"url"`
	CreatedAt time.Time `bson:"created_at"`
}

type ChapterComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	ChapterID primitive.ObjectID `bson:"chapter_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Comment   string             `bson:"comment"`
	CreatedAt time.Time          `bson:"created_at"`
}

type Volume struct {
	ID       primitive.ObjectID   `bson:"_id"`
	MangaID  primitive.ObjectID   `bson:"manga_id"`
	Chapters []primitive.ObjectID `bson:"chapters"`
}

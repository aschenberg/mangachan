// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package pgdb

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CasbinRule struct {
	ID    int32
	Ptype pgtype.Text
	V0    pgtype.Text
	V1    pgtype.Text
	V2    pgtype.Text
	V3    pgtype.Text
	V4    pgtype.Text
	V5    pgtype.Text
}

type Genre struct {
	GenreID   int64
	Title     string
	CreatedAt int64
	UpdatedAt int64
}

type Manga struct {
	MangaID   int64
	Title     string
	Titles    []string
	Synonyms  []string
	CoverID   int64
	Type      string
	Country   string
	Status    pgtype.Text
	CreatedAt int64
	UpdatedAt int64
}

type MangaCover struct {
	CoverID     int64
	CoverDetail pgtype.Text
	Thumbnail   pgtype.Text
	Extra       []string
	CreatedAt   int64
	UpdatedAt   int64
}

type MangaDetail struct {
	DetailID  int64
	Published pgtype.Text
	Authors   []string
	Artist    []string
	Source    pgtype.Text
	Summary   pgtype.Text
	UpdatedAt int64
	CreatedAt int64
}

type MangaGenre struct {
	MgID      int64
	MangaID   int64
	GenreID   int64
	CreatedAt int64
	UpdatedAt int64
}

type MangaScore struct {
	ScoreID   int64
	Score     pgtype.Numeric
	CreatedAt int64
	UpdatedAt int64
}

type Setting struct {
	SettingID    int32
	FirstVoucher bool
}

type User struct {
	UserID       int64
	AppID        string
	Email        string
	Picture      pgtype.Text
	Role         int16
	IsActive     bool
	GivenName    pgtype.Text
	FamilyName   pgtype.Text
	Name         pgtype.Text
	RefreshToken string
	IsDeleted    bool
	CreatedAt    int64
	UpdatedAt    int64
}

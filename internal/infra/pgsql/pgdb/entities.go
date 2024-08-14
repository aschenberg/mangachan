// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package pgdb

import (
	"github.com/jackc/pgx/v5/pgtype"
)

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

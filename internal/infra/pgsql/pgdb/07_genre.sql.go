// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: 07_genre.sql

package pgdb

import (
	"context"
)

const createGenre = `-- name: CreateGenre :exec
INSERT INTO genre (
  genre_id,
  title,
  updated_at,
  created_at
) VALUES (
  $1,$2,$3,$4
) ON CONFLICT (title) DO NOTHING
`

type CreateGenreParams struct {
	GenreID   int64
	Title     string
	UpdatedAt int64
	CreatedAt int64
}

func (q *Queries) CreateGenre(ctx context.Context, arg CreateGenreParams) error {
	_, err := q.db.Exec(ctx, createGenre,
		arg.GenreID,
		arg.Title,
		arg.UpdatedAt,
		arg.CreatedAt,
	)
	return err
}

const findGenreByTitle = `-- name: FindGenreByTitle :many
SELECT genre_id FROM genre
WHERE title = ANY($1::text[])
`

func (q *Queries) FindGenreByTitle(ctx context.Context, titles []string) ([]int64, error) {
	rows, err := q.db.Query(ctx, findGenreByTitle, titles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var genre_id int64
		if err := rows.Scan(&genre_id); err != nil {
			return nil, err
		}
		items = append(items, genre_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

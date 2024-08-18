-- name: CreateGenre :exec
INSERT INTO genre (
  genre_id,
  title,
  updated_at,
  created_at
) VALUES (
  $1,$2,$3,$4
) ON CONFLICT (title) DO NOTHING;
-- name: FindGenreByTitle :many
SELECT genre_id FROM genre
WHERE title = ANY(@titles::text[]);




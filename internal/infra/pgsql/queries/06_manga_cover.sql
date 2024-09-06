-- name: CreateMangaCover :one
INSERT INTO manga_cover (
  cover_id,
  cover_detail,
  thumbnail,
  extra,
  updated_at,
  created_at
) VALUES (
  $1,$2,$3,$4,$5,$6
) RETURNING cover_id;




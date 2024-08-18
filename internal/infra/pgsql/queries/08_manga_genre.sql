-- name: CreateMangaGenre :exec
INSERT INTO manga_genre (
  mg_id,
  manga_id,
  genre_id,
  updated_at,
  created_at
) VALUES (
  $1,$2,$3,$4,$5
);




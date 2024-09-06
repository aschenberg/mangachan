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


-- name: FindMangaGenre :many
SELECT * FROM manga_genre AS a 
LEFT JOIN genre AS b ON a.genre_id = b.genre_id 
WHERE a.manga_id = $1;




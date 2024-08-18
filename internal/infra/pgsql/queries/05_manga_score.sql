-- name: CreateMangaScore :exec
INSERT INTO manga_score (
  score_id,
	score,
	updated_at,
	created_at
) VALUES (
  $1,$2,$3,$4
);




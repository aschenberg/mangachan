-- name: CreateManga :one
INSERT INTO manga (
  	manga_id,
	title,
	title_en,
  	synonyms,
	cover_id,
	type,
	country,
	status,
	updated_at,
	created_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
) ON CONFLICT (title) DO NOTHING 
RETURNING manga_id,title;

-- name: FindMangaByID :one
SELECT * FROM manga
WHERE manga_id = $1 LIMIT 1;


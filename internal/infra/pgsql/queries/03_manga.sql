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
SELECT * FROM manga AS a 
LEFT JOIN manga_detail AS b ON a.manga_id = b.detail_id 
LEFT JOIN manga_score AS c ON a.manga_id = c.score_id
LEFT JOIN manga_cover AS d ON a.manga_id = d.cover_id
WHERE a.manga_id = $1 LIMIT 1;


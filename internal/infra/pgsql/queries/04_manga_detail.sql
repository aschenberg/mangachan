-- name: CreateMangaDetail :exec
INSERT INTO manga_detail (
  detail_id,
	published,
	authors,
    artist,
	summary,
	updated_at,
	created_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7
);




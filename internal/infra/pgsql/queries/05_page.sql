-- name: InsertPage :one
INSERT INTO page (
	lesson_id,
	page_number,
	created_at,
	updated_at
) VALUES (
  $1,$2,$3,$4
)
RETURNING page_id,lesson_id,page_number;


-- name: FindPageByID :one
SELECT * FROM page
WHERE page_id = $1 LIMIT 1;


-- name: RemovePage :exec
DELETE FROM page
WHERE page_id = $1;


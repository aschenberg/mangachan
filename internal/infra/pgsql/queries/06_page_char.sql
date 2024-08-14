-- name: InsertPageChar :one
INSERT INTO page_char (
	page_id,
	char_id,
	created_at,
	updated_at
) VALUES (
  $1,$2,$3,$4
)
RETURNING page_char_id,char_id;


-- name: FindPageCharByID :one
SELECT * FROM page_char
WHERE page_char_id = $1 LIMIT 1;


-- name: RemovePageChar :exec
DELETE FROM page_char
WHERE page_char_id = $1;


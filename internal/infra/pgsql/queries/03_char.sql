-- name: InsertChar :one
INSERT INTO char (
	step,
	image,
	voice,
  	en,
	id,
	is_deleted,
	created_at,
	updated_at
) VALUES (
  $1,$2,$3,$4,$5, $6,$7,$8
)
RETURNING char_id,image;


-- name: FindCharByID :one
SELECT * FROM char
WHERE char_id = $1 LIMIT 1;


-- name: RemoveChar :exec
DELETE FROM char
WHERE char_id = $1;


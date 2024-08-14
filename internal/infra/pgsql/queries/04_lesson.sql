-- name: InsertLesson :one
INSERT INTO lesson (
	title,
	video_url,
	number,
  	is_free,
	is_deleted,
	created_at,
	updated_at
) VALUES (
  $1,$2,$3,$4,$5, $6,$7
)
RETURNING lesson_id,title;


-- name: FindLessonByID :one
SELECT * FROM lesson
WHERE lesson_id = $1 LIMIT 1;


-- name: RemoveLesson :exec
DELETE FROM lesson
WHERE lesson_id = $1;


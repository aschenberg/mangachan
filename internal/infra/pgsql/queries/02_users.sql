-- name: FindBySubID :one
SELECT * FROM users
WHERE sub_id = $1 LIMIT 1;



-- name: CreateUser :one
INSERT INTO users (
  user_id,
	app_id,
	sub_id,
	email,
  picture,
	role,
	is_active,
	given_name,
	family_name,
	refresh_token,
	created_at,
	is_deleted,
	updated_at
) VALUES (
  $1,$2,$3,$4,$5, $6,$7,$8,$9,$10,$11,$12,$13
)
RETURNING user_id,created_at,updated_at;

-- name: GetRefreshToken :one
SELECT refresh_token FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;



-- name: UpdateUser :exec
UPDATE users SET 
email = COALESCE(sqlc.narg('email'),email)  
WHERE user_id = @user_id;

-- name: GetUserActiveStatus :one
SELECT is_active FROM users
WHERE email = $1 LIMIT 1;

-- name: CheckUserEmailExist :one
SELECT email FROM users
WHERE email = $1 LIMIT 1;


-- name: UpdateUserToken :exec
UPDATE users SET refresh_token = $2 WHERE user_id = $1;

-- name: UpdateUserActive :exec
UPDATE users SET is_active = $1 WHERE email = $2;


-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;


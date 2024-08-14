-- name: FindBySubID :one
SELECT * FROM users
WHERE app_id = $1 LIMIT 1;



-- name: CreateOrUpdateUser :one
INSERT INTO users (
  user_id,
	app_id,
	email,
  picture,
	role,
	is_active,
	given_name,
	family_name,
	name,
	refresh_token,
	is_deleted,
	created_at,
	updated_at
) VALUES (
  @id,@appId,@email,@picture,@role, @isActive,@givenName,@familyName,@name,@refreshToken,@isDeleted,@createdAt,@updatedAt
) ON CONFLICT (app_id) DO UPDATE SET 
picture = COALESCE(sqlc.narg('picture'),picture),
given_name = COALESCE(sqlc.narg('given_name'),given_name),
family_name = COALESCE(sqlc.narg('family_name'),family_name),
name = COALESCE(sqlc.narg('name'),name)
RETURNING *, CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation;

-- name: UpdateRefreshToken :exec
UPDATE users SET 
refresh_token = @refreshToken  
WHERE app_id = @app_id;

-- name: GetRefreshToken :one
SELECT refresh_token FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;





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


CREATE TABLE IF NOT EXISTS users (
	user_id BIGINT PRIMARY KEY,
	app_id VARCHAR(100) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	picture TEXT,
	role SMALLINT NOT NULL,
	is_active BOOLEAN NOT NULL,
	given_name TEXT,
	family_name TEXT,
	name TEXT,
	refresh_token TEXT NOT NULL,
	is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL
);

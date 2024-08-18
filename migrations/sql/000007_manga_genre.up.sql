CREATE TABLE IF NOT EXISTS manga_genre (
	mg_id BIGINT PRIMARY KEY,
    manga_id BIGINT NOT NULL,
    genre_id BIGINT NOT NULL,
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL
);
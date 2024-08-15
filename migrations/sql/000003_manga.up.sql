CREATE TABLE IF NOT EXISTS manga (
	manga_id BIGINT PRIMARY KEY,
	title VARCHAR NOT NULL,
	title_en VARCHAR,
	synonyms TEXT,
	cover_id BIGINT NOT NULL,
	type VARCHAR NOT NULL,
	country VARCHAR NOT NULL,
	status VARCHAR,
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL
);
CREATE TABLE IF NOT EXISTS manga_score (
	score_id BIGINT PRIMARY KEY,
	score NUMERIC(4,2),
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL
);
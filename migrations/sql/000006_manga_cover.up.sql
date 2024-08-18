CREATE TABLE IF NOT EXISTS manga_cover (
	cover_id BIGINT PRIMARY KEY,
	cover_detail TEXT,
    thumbnail TEXT,
    extra TEXT[],
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL
);
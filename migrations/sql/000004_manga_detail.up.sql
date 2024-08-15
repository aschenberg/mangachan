CREATE TABLE IF NOT EXISTS manga_detail (
	detail_id BIGINT PRIMARY KEY,
	published VARCHAR,
	authors TEXT[],
	artist TEXT[],
	summary TEXT,
    updated_at BIGINT NOT NULL,
	created_at BIGINT NOT NULL
	
);
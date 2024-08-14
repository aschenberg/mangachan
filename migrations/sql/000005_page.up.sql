CREATE TABLE IF NOT EXISTS page (
	page_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	lesson_id INT NOT NULL,
    page_number INT NOT NULL,
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL,
	CONSTRAINT page_lesson_id_fkey FOREIGN KEY (lesson_id) REFERENCES lesson (lesson_id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS page_char (
	page_char_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	page_id INT NOT NULL,
    char_id INT NOT NULL,
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL,
	CONSTRAINT page_char_page_id_fkey FOREIGN KEY (page_id) REFERENCES page (page_id) ON DELETE CASCADE,
    CONSTRAINT page_char_char_id_fkey FOREIGN KEY (char_id) REFERENCES char (char_id) ON DELETE CASCADE
);
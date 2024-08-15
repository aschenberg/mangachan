package entity

import "time"

type Manga struct {
	Manga_ID     int64     `json:"manga_id"`
	Title        string    `json:"title"`
	TitleEnglish string    `json:"title_en"`
	Synonyms     []string  `json:"synonyms"`
	CoverID      int64     `json:"cover_id"`
	Type         string    `json:"type"`
	Country      string    `json:"country"`
	Published    string    `json:"published"`
	Status       string    `json:"status"`
	Authors      []string  `json:"authors"`
	Artists      []string  `json:"artists"`
	Score        float64   `json:"score"`
	Summary      string    `json:"summary"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

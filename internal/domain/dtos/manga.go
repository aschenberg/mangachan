package dtos

import "manga/internal/domain/models"

type CreateManga struct {
	Title        string            `json:"title"`
	TitleEnglish string            `json:"title_en"`
	Synonyms     []string          `json:"synonyms"`
	Type         string            `json:"type"`
	Country      string            `json:"country"`
	Cover        models.MangaCover `json:"cover"`
	Published    string            `json:"published"`
	Status       string            `json:"status"`
	Authors      []string          `json:"authors"`
	Artists      []string          `json:"artists"`
	Genres       []string          `json:"genres"`
	Score        float64           `json:"score"`
	Themes       []string          `json:"themes"`
	Demographic  []string          `json:"demographic"`
	Summary      string            `json:"summary"`
}

type CreatedMangaResponse struct {
	MangaID int64
	Title   string
}

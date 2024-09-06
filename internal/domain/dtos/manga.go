package dtos

import (
	"fmt"
	"manga/internal/domain/models"
	"manga/internal/infra/pgsql/pgdb"
	"manga/pkg/utils"
)

type CreateManga struct {
	Source       string            `json:"source"`
	Method       string            `json:"method"`
	Request      any               `json:"request"`
	MangaID      string            `json:"manga_id"`
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

type Manga struct {
	Source       string            `json:"source"`
	MangaID      string            `json:"manga_id"`
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
	Genres       []models.Genre    `json:"genres"`
	Score        float64           `json:"score"`
	Summary      string            `json:"summary"`
}

func ToManga(v pgdb.FindMangaByIDRow, w []pgdb.FindMangaGenreRow) Manga {
	var a []models.Genre
	for _, v := range w {
		fmt.Print(v.Title.String)
		a = append(a, models.Genre{
			GenreID: utils.Int64ToStr(v.GenreID),
			Name:    v.Title.String})
	}
	return Manga{
		Source:       "",
		MangaID:      utils.Int64ToStr(v.MangaID),
		Title:        v.Title,
		TitleEnglish: v.TitleEn.String,
		Synonyms:     v.Synonyms,
		Type:         v.Type,
		Country:      v.Country,
		Published:    v.Published.String,
		Status:       v.Status.String,
		Artists:      v.Artist,
		Score:        utils.PgNumToFloat(v.Score),
		Authors:      v.Authors,
		Cover: models.MangaCover{
			CoverID:     utils.Int64ToStr(v.CoverID_2.Int64),
			CoverDetail: v.CoverDetail.String,
			Thumbnail:   v.Thumbnail.String,
			Extra:       v.Extra},
		Summary: v.Summary.String,
		Genres:  a,
	}
}

type IndexedManga struct{
	ID    		 string `json:"id"`
	Title        string            `json:"title"`
	TitleEnglish string            `json:"title_en"`
	Synonyms     []string          `json:"synonyms"`
	Type         string            `json:"type"`
	Cover        string             `json:"cover"`
	Status       string            `json:"status"`
	Genres       []models.Genre          `json:"genres"`
	Score        float64           `json:"score"`

}

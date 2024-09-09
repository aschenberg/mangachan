package dtos

import (
	"manga/internal/domain/models"
	"manga/internal/infra/pgsql/pgdb"
	"manga/pkg/utils"
	"reflect"
)

type CreateManga struct {
	Source       string            `json:"source"`
	Method       string            `json:"method"`
	Request      any               `json:"request"`
	MangaID      string            `json:"manga_id"`
	Title        string            `json:"title"`
	Titles []string            `json:"titles"`
	Synonyms     []string          `json:"synonyms"`
	Type         string            `json:"type"`
	Country      string            `json:"country"`
	Cover        models.MangaCover `json:"cover"`
	Published    string            `json:"published"`
	Status       string            `json:"status"`
	Authors      []string          `json:"authors"`
	Artists      []string          `json:"artists"`
	Genres       []models.Genre          `json:"genres"`
	Score        float64           `json:"score"`
	Themes       []string          `json:"themes"`
	Demographic  []string          `json:"demographic"`
	Summary      string            `json:"summary"`
}

type CreatedMangaResponse struct {
	MangaID string
	Title   string
}

type Manga struct {
	Source       string            `json:"source"`
	MangaID      string            `json:"manga_id"`
	Title        string            `json:"title"`
	Titles []string            `json:"titles"`
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
		a = append(a, models.Genre{
			GenreID: utils.Int64ToStr(v.GenreID),
			Name:    v.Title.String})
	}
	return Manga{
		Source:       "",
		MangaID:      utils.Int64ToStr(v.MangaID),
		Title:        v.Title,
		Titles: v.Titles,
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
	Titles []string            `json:"titles"`
	Synonyms     []string          `json:"synonyms"`
	Type         string            `json:"type"`
	Cover        string             `json:"cover"`
	Status       string            `json:"status"`
	Genres       []models.Genre          `json:"genres"`
	Score        float64           `json:"score"`

}

func CreateByMangaUpdate(v models.MangaUpdateByID )CreateManga{
	return CreateManga{
		Title:        v.Title,
			Titles: MguMap(v.Associated, "Title"),
			Synonyms:    []string{},
			Type:         v.Type,
			Published:    v.Year,
			Country:      switchCountry(v.Type),
			Status:       v.Status,
			Authors:      MguMap(v.Authors, "Name"),
			Artists:      MguMap(v.Authors, "Name"),
			Genres:       ToMangaGenre(v.Genres, "Genre"),
			Themes:       MguMap(v.Categories, "Category"),
			Demographic:  []string{},
			Summary:      v.Description,
			Score:        float64(v.BayesianRating),
			Cover: models.MangaCover{
				CoverDetail: v.Image.URL.Original,
				Thumbnail:   v.Image.URL.Original,
			},
	}
} 

func CreateByMyAnimeList(v models.JikanMangaByID)CreateManga{
	return CreateManga{Title:        v.Data.Title,
		Titles: MguMap(v.Data.Titles,"Title"),
		Synonyms:     v.Data.TitleSynonyms,
		Type:         v.Data.Type,
		Published:    v.Data.Published.String,
		Country:      switchCountry(v.Data.Type),
		Status:       v.Data.Status,
		Authors:      addName(v.Data.Authors),
		Artists:      addName(v.Data.Authors),
		Genres:       ToMangaGenre(v.Data.Genres,"Name"),
		Themes:       addName(v.Data.Themes),
		Demographic:  addName(v.Data.Demographics),
		Summary:      v.Data.Synopsis,
		Score:        v.Data.Score,
		Cover: models.MangaCover{
			CoverDetail: v.Data.Images.Webp.LargeImageURL,
			Thumbnail:   v.Data.Images.Webp.ImageURL,
		},}
}


func addName(vals []models.MALSubJson) []string {
	var valStr []string
	for _, a := range vals {
		valStr = append(valStr, a.Name)
	}
	return valStr
}

func MguMap(vals interface{}, field string) []string {
	valsValue := reflect.ValueOf(vals)
	if valsValue.Kind() != reflect.Slice {
		panic("mguMap: input must be a slice")
	}

	var valStr []string
	for i := 0; i < valsValue.Len(); i++ {
		elem := valsValue.Index(i)
		field := elem.FieldByName(field)
		if !field.IsValid() {
			panic("mguMap: struct does not have a Title field")
		}
		valStr = append(valStr, field.String())
	}
	return valStr
}
func ToMangaGenre(vals interface{}, field string)[]models.Genre{
    valsValue := reflect.ValueOf(vals)
	if valsValue.Kind() != reflect.Slice {
		panic("mguMap: input must be a slice")
	}

	var valStr []models.Genre
	for i := 0; i < valsValue.Len(); i++ {
		elem := valsValue.Index(i)
		field := elem.FieldByName(field)
		if !field.IsValid() {
			panic("mguMap: struct does not have a Title field")
		}
		valStr = append(valStr, models.Genre{Name: field.String()})
	}
	return valStr
}
func ToGenreFromPGRB(val []pgdb.FindMangaGenreRow )[]models.Genre{

	var valStr []models.Genre
	for _, v := range val {
		valStr = append(valStr, models.Genre{GenreID: 
			utils.Int64ToStr(v.GenreID),
		Name: v.Title.String,})
	}
	return valStr
}

func switchCountry(val string) string {
	switch val {
	case "Manga":
		return "jp"
	case "manga":
		return "jp"
	case "Manhwa":
		return "kr"
	case "manhwa":
		return "kr"
	case "Manhua":
		return "cn"
	case "manhua":
		return "cn"
	default:
		return "unknown"
	}
}

package models

import "time"

type JikanMangaByID struct {
	Data struct {
		MalID  int    `json:"mal_id"`
		URL    string `json:"url"`
		Images struct {
			Jpg struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"jpg"`
			Webp struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"webp"`
		} `json:"images"`
		Approved bool `json:"approved"`
		Titles   []struct {
			Type  string `json:"type"`
			Title string `json:"title"`
		} `json:"titles"`
		Title         string      `json:"title"`
		TitleEnglish  string      `json:"title_english"`
		TitleJapanese string      `json:"title_japanese"`
		TitleSynonyms []string    `json:"title_synonyms"`
		Type          string      `json:"type"`
		Chapters      interface{} `json:"chapters"`
		Volumes       int         `json:"volumes"`
		Status        string      `json:"status"`
		Publishing    bool        `json:"publishing"`
		Published     struct {
			From time.Time `json:"from"`
			To   time.Time `json:"to"`
			Prop struct {
				From struct {
					Day   int `json:"day"`
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"from"`
				To struct {
					Day   int `json:"day"`
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"to"`
			} `json:"prop"`
			String string `json:"string"`
		} `json:"published"`
		Score          float64       `json:"score"`
		Scored         float64       `json:"scored"`
		ScoredBy       int           `json:"scored_by"`
		Rank           int           `json:"rank"`
		Popularity     int           `json:"popularity"`
		Members        int           `json:"members"`
		Favorites      int           `json:"favorites"`
		Synopsis       string        `json:"synopsis"`
		Background     string        `json:"background"`
		Authors        []MALSubJson  `json:"authors"`
		Serializations []MALSubJson  `json:"serializations"`
		Genres         []MALSubJson  `json:"genres"`
		ExplicitGenres []interface{} `json:"explicit_genres"`
		Themes         []MALSubJson  `json:"themes"`
		Demographics   []MALSubJson  `json:"demographics"`
	} `json:"data"`
}

type MALSubJson struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

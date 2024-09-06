package models

type MangaUpdateByID struct {
	SeriesID    int            `json:"series_id"`
	Title       string         `json:"title"`
	URL         string         `json:"url"`
	Associated  []MguAssosiate `json:"associated"`
	Description string         `json:"description"`
	Image       struct {
		URL struct {
			Original string `json:"original"`
			Thumb    string `json:"thumb"`
		} `json:"url"`
	} `json:"image"`
	Type           string  `json:"type"`
	Year           string  `json:"year"`
	BayesianRating float64 `json:"bayesian_rating"`
	Genres         []struct {
		Genre string `json:"genre"`
	} `json:"genres"`
	Categories []struct {
		SeriesID   int    `json:"series_id"`
		Category   string `json:"category"`
		Votes      int    `json:"votes"`
		VotesPlus  int    `json:"votes_plus"`
		VotesMinus int    `json:"votes_minus"`
		AddedBy    int    `json:"added_by"`
	} `json:"categories"`
	Status string `json:"status"`

	Authors []struct {
		Name     string `json:"name"`
		AuthorID int    `json:"author_id"`
		Type     string `json:"type"`
	} `json:"authors"`
	Publishers []struct {
		PublisherName string `json:"publisher_name"`
		PublisherID   int    `json:"publisher_id"`
		Type          string `json:"type"`
		Notes         string `json:"notes"`
	} `json:"publishers"`
	Publications []struct {
		PublicationName string `json:"publication_name"`
		PublisherName   string `json:"publisher_name"`
		PublisherID     int    `json:"publisher_id"`
	} `json:"publications"`
}

type MguAssosiate struct {
	Title string `json:"title"`
}

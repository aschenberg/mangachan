package enum

type Source int16

const (
	MyAnimeList Source = iota + 1
	MangaUpdate
)

func (s Source) String() string {
	return [...]string{"Generation One", "Generation Two"}[s]
}

func (s Source) ValueToString(val int16) string {
	return [...]string{"", "Generation One", "Generation Two"}[val]
}

type SourceListResponse struct {
	Name  string `json:"name"`
	Value Source `json:"value"`
	Url   string `json:"url"`
}

var SourceName = map[Source]SourceInfo{
	MyAnimeList: {Name: "My Anime List", Source: MyAnimeList, Url: "https://api.jikan.moe/v4"},
	MangaUpdate: {Name: "Manga Update", Source: MangaUpdate, Url: "https://api.mangaupdates.com"},
}

type SourceInfo struct {
	Name   string
	Source Source
	Url    string
}

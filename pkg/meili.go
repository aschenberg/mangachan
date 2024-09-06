package pkg

import (
	"fmt"
	"manga/config"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearch struct {
	client *meilisearch.ServiceManager
	index  string
}

func NewMeili(cfg *config.Config) meilisearch.ServiceManager {
	connectionString := fmt.Sprintf("%s:%s", cfg.Meili.Host, cfg.Meili.Port)
	client := meilisearch.New(connectionString, meilisearch.WithAPIKey(cfg.Meili.ApiKey))
	return client
}

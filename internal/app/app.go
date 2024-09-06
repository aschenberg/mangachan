package app

import (
	"context"
	"fmt"
	"manga/api/route"
	"manga/config"
	"manga/internal/domain/dtos"
	pgsql "manga/internal/infra/pgsql/repository"
	"manga/pkg"
	"manga/pkg/logging"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func InitServer(cfg *config.Config, log logging.Logger) *pkg.Server {

	// arangoC := arango.NewArangoDatabase(cfg)
	// db := arango.ConnectDatabase(arangoC, cfg)
	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!!!")
	})
	// OpenID Initialized
	openIDProvider := pkg.NewOidcProvider(cfg)
	openIDCfg := pkg.NewOidcConfig(cfg, openIDProvider)

	ml := pkg.NewMeili(cfg)

	// Redis Initialized
	redisC, err := pkg.NewRedis(cfg)
	if err != nil {
		log.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	// S3 Storage Initialized
	minio, err := pkg.NewS3Storage(cfg, log)
	if err != nil {
		log.Fatal(logging.S3, logging.Startup, err.Error(), nil)
	}

	// PostgreSQL Initialized
	pg, err := pkg.NewPg(cfg)
	if err != nil {
		log.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	// Rabbit Mq Initialized
	rmq, err := pkg.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal(logging.Rabbit, logging.Startup, err.Error(), nil)
	}

	defer rmq.Close()

	err = rmq.Consume("manga", func(d amqp.Delivery) {
		if string(d.Body) != "" {
			mgr := pgsql.NewMangaRepository(pg, log)

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			mg, err := mgr.FindbyId(ctx, string(d.Body))
			if err != nil {
				fmt.Sprintf(err.Error())
			}
			genre, err := mgr.FindMangaGenre(ctx, string(d.Body))
			if err != nil {
				fmt.Sprintf(err.Error())
			}
			res := dtos.ToManga(mg, genre)

			index := ml.Index("Manga")
			documents := []map[string]interface{}{
				{"id": res.MangaID,
					"title":    res.Title,
					"thumb":    res.Cover.Thumbnail,
					"synonyms": res.Synonyms,
					"genres":   res.Genres,
					"score":    res.Score,
					"status":   res.Status,
					"type":     res.Type},
			}
			_, err = index.AddDocuments(documents)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	})

	// Gin Route Initialized
	api := r.Group("/api")
	v1 := api.Group("/v1")
	m := v1.Group("/m")
	{
		auth := m.Group("/auth")
		route.Auth(auth, cfg, pg, openIDCfg, openIDProvider, log)
		manga := m.Group("/manga")
		route.Manga(manga, cfg, pg, redisC, minio, rmq, ml, log)
	}

	httpserverServer := pkg.NewHttp(cfg, r)
	return httpserverServer
}

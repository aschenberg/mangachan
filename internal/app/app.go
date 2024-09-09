package app

import (
	"manga/api/route"
	"manga/config"
	"manga/internal/infra/rabbitmq"
	"manga/pkg"
	"manga/pkg/logging"

	"github.com/gin-gonic/gin"
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
	rmq , err := pkg.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	
	rmqtask, err := rabbitmq.NewTask(rmq.Channel)
	if err != nil {
		log.Fatal(logging.Rabbit, logging.Startup, err.Error(), nil)
	}

	// Gin Route Initialized
	api := r.Group("/api")
	v1 := api.Group("/v1")
	m := v1.Group("/m")
	{
		auth := m.Group("/auth")
		route.Auth(auth, cfg, pg, openIDCfg, openIDProvider, log)
		manga := m.Group("/manga")
		route.Manga(manga, cfg, pg, redisC, minio, rmqtask, ml, log)
	}

	httpserverServer := pkg.NewHttp(cfg, r)
	return httpserverServer
}

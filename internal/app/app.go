package app

import (
	"manga/api/route"
	"manga/config"
	"manga/pkg/cache"
	"manga/pkg/httpserver"
	"manga/pkg/logging"
	"manga/pkg/openid"
	"manga/pkg/postgres"
	"manga/pkg/s3"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

func InitServer(cfg *config.Config, log logging.Logger) *httpserver.Server {

	// arangoC := arango.NewArangoDatabase(cfg)
	// db := arango.ConnectDatabase(arangoC, cfg)

	openIDProvider := openid.NewOidcProvider(cfg)
	openIDCfg := openid.NewOidcConfig(cfg, openIDProvider)
	redisC, err := cache.InitRedis(cfg)
	if err != nil {
		log.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}
	minio, err := s3.InitS3Storage(cfg, log)
	if err != nil {
		log.Fatal(logging.S3, logging.Startup, err.Error(), nil)
	}
	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!!!")
	})
	pg, err := postgres.InitPg(cfg)
	if err != nil {
		log.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	RegisterRoutes(r, cfg, pg, openIDCfg, openIDProvider, redisC, minio, log)

	httpserverServer := httpserver.New(cfg, r)
	return httpserverServer
}

func RegisterRoutes(gin *gin.Engine, cfg *config.Config, pg *postgres.Postgres, oidcCfg *oauth2.Config, oidcPvd *oidc.Provider, redis *redis.Client, minio *minio.Client, log logging.Logger) {

	// Al l Public APIs
	api := gin.Group("/api")
	v1 := api.Group("/v1")
	m := v1.Group("/m")
	{
		auth := m.Group("/auth")
		route.Auth(auth, cfg, pg, oidcCfg, oidcPvd, log)
		manga := m.Group("/manga")
		route.Manga(manga, cfg, pg, redis, minio, log)
	}

}

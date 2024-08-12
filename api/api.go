package api

import (
	"manga/api/route"
	"manga/config"
	"manga/pkg/arango"
	"manga/pkg/httpserver"
	"manga/pkg/openid"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func InitServer(cfg *config.Config) *httpserver.Server {

	arangoC := arango.NewArangoDatabase(cfg)
	db := arango.ConnectDatabase(arangoC, cfg)
	openIDProvider := openid.NewOidcProvider(cfg)
	openIDCfg := openid.NewOidcConfig(cfg, openIDProvider)

	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!!!")
	})

	RegisterRoutes(r, cfg, db, openIDCfg, openIDProvider)

	httpserverServer := httpserver.New(cfg, r)
	return httpserverServer
}

func RegisterRoutes(gin *gin.Engine, cfg *config.Config, db arangodb.Database, oidcCfg *oauth2.Config, oidcPvd *oidc.Provider) {

	// Al l Public APIs
	api := gin.Group("/api")
	v1 := api.Group("/v1")
	{
		auth := v1.Group("/auth")
		route.Auth(auth, cfg, db, oidcCfg, oidcPvd)
	}

}

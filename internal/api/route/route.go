package route

import (
	"manga/config"
	"manga/db"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

func Setup(env *config.Env, timeout time.Duration, db db.MongoDB, gin *gin.Engine, oidc *oauth2.Config, provider *oidc.Provider, redis *redis.Client) {

	// Al l Public APIs
	publicV1 := gin.Group("/api/v1/pub")
	NewAuthRouter(env, timeout, db, publicV1, oidc, provider)
	NewMangaRouter(env, timeout, db, publicV1, redis)
	// NewLoginRouter(env, timeout, db, publicRouter)
	// NewRefreshTokenRouter(env, timeout, db, publicRouter)
	// NewSignoutRouter(publicRouter)
	// privateV1 := gin.Group("/api/v1")
	// protectedRouter := gin.Group("/api/v1")
	// Middleware to verify AccessToken
	// protectedRouter.Use(middleware.SupertokensMiddleware)
	// All Private APIs
	// NewProfileRouter(env, timeout, db, protectedRouter)

	// NewHospitalRouter(env, timeout, db, protectedRouter)
	// NewUserRouter(env, timeout, db, protectedRouter)
	// NewWorkbookRouter(env, timeout, db, protectedRouter)
}

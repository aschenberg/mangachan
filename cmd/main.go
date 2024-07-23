package main

import (
	"context"
	"fmt"
	"manga/api/route"
	"manga/config"
	"manga/pkg/cache"
	"manga/pkg/mongodb"
	"manga/pkg/oidc"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {

	r := gin.Default()

	cfg := config.GetConfig()

	db := mongodb.NewMongoDatabase(cfg)
	mongodb.CloseMongoDBConnection(db)
	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}
	// redis := config.NewRedis(redisClient)

	oidcProvider := oidc.NewOidcProvider(cfg)
	oidcCon := oidc.NewOidcConfig(cfg, oidcProvider)

	client := cache.GetRedis()
    pubsub := client.PSubscribe(context.Background(), "__keyevent@0__:expired")
	// Ensure subscription is closed on exit
	defer pubsub.Close()

	// Listen for messages
	go func() {
		for msg := range pubsub.Channel() {
			handleExpirationEvent(msg)
		}
	}()

	// config := cors.DefaultConfig()

	// config.AllowOrigins = []string{"http://localhost:5555", "http://51.79.185.128:5555", "http://localhost:5173"}
	// config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "PATCH", "DELETE"}
	// config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	// config.ExposeHeaders = []string{"Content-Length"}
	// config.AllowPrivateNetwork = true
	// config.AllowCredentials = true
	// config.MaxAge = 12 * time.Hour
	// r.Use(cors.New(config))
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!!!")
	})
	// r.GET("/image/:imageId", func(c *gin.Context) {
	// 	serveImage(c, db)
	// })
	// CORS

	route.Setup(env, timeout, db, r, oidcCon, oidcProvider, redisClient)

	r.Run(env.ServerAddress)

}

func handleExpirationEvent(msg *redis.Message) {
	fmt.Printf("Received expiration event for key: %s\n", msg.Payload)
}

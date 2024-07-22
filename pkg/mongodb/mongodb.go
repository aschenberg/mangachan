package mongodb

import (
	"context"
	"fmt"
	"log"
	"manga/config"
	"manga/db"
	"time"
)

func NewMongoDatabase(cfg *config.Config) db.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbURI := fmt.Sprintf("mongodb://%s:%s/%s?authSource=admin", cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.DbName)

	client, err := db.NewClient(ctx, mongodbURI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx); err != nil {
		panic(err)
	}

	return client
}

func CloseMongoDBConnection(client db.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}

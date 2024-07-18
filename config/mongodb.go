package config

import (
	"context"
	"fmt"
	"log"
	"manga/db"
	"time"
)

func NewMongoDatabase(env *Env) db.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	isDbAtlas := env.MongoAtlas
	dbUrl := env.DBUrl
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName

	mongodbURI := fmt.Sprintf("mongodb://%s:%s/%s?authSource=admin", dbHost, dbPort, dbName)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}
	if isDbAtlas {
		mongodbURI = dbUrl
	}
	client, err := db.NewClient(ctx, mongodbURI, dbUser, dbPass)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx); err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

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

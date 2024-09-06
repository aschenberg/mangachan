package pkg

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"manga/config"
	"net"
	"net/http"
	"time"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"github.com/arangodb/go-driver/v2/connection"
)

func NewArangoDatabase(cfg *config.Config) arangodb.Client {
	var err error
	// Create an HTTP connection to the database
	singleConnectionUrl := fmt.Sprintf("http://%s:%s", cfg.Arango.Host, cfg.Arango.Port)
	endpoint := connection.NewRoundRobinEndpoints([]string{singleConnectionUrl})
	// conn, err := connection.NewPool(10, nil)
	// if err != nil {

	// }
	authRoot := connection.NewBasicAuth("root", cfg.Arango.Password)

	rootClient := arangodb.NewClient(
		connection.NewHttpConnection(
			connection.HttpConfiguration{
				Endpoint:       endpoint,
				Authentication: authRoot,
				ContentType:    connection.ApplicationJSON}))

	err = createUser(rootClient, cfg.Arango.User, cfg.Arango.Password)
	if err != nil {
		if shared.IsArangoErrorWithErrorNum(err, shared.ErrUserDuplicate) {
			fmt.Println("User already exists, skipping creation.")
		} else {
			log.Fatalf("Failed to create user: %v", err)
		}
	} else {
		fmt.Printf("User '%s' created successfully.\n", cfg.Arango.User)
	}
	// Check if the database exists
	exists, err := databaseExists(rootClient, cfg.Arango.DbName)
	if err != nil {
		log.Fatalf("Failed to check if database exists: %v", err)
	}
	if !exists {
		// Create the database if it doesn't exist
		if err := createDatabase(rootClient, cfg.Arango.DbName, cfg); err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Database '%s' created.\n", cfg.Arango.DbName)
	} else {
		fmt.Printf("Database '%s' already exists.\n", cfg.Arango.DbName)
	}

	// Now authenticate with the new users
	authUser := connection.NewBasicAuth(cfg.Arango.User, cfg.Arango.Password)
	newClient := arangodb.NewClient(
		connection.NewHttpConnection(
			connection.HttpConfiguration{
				Authentication: authUser, Endpoint: endpoint, ContentType: connection.ApplicationJSON}))

	return newClient
}

// Check if a database exists
func databaseExists(client arangodb.Client, dbName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database, err := client.DatabaseExists(ctx, dbName)
	if err != nil {

		return false, err
	}

	return database, nil
}

// Create a database
func createDatabase(client arangodb.Client, dbName string, cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.CreateDatabase(ctx, dbName, &arangodb.CreateDatabaseOptions{Users: []arangodb.CreateDatabaseUserOptions{{UserName: cfg.Arango.User, Password: cfg.Arango.Password}}})
	return err
}

func jSONHTTPConnectionConfig(endpoint connection.Endpoint) connection.HttpConfiguration {
	return connection.HttpConfiguration{
		Endpoint:       endpoint,
		ContentType:    connection.ApplicationJSON,
		ArangoDBConfig: connection.ArangoDBConfiguration{},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 90 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// Check if a user exists
func userExists(client arangodb.Client, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.UserExists(ctx, username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Create a user
func createUser(client arangodb.Client, username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.CreateUser(ctx, username, &arangodb.UserOptions{
		Password: password,
	})
	return err
}

func ConnectDatabase(client arangodb.Client, cfg *config.Config) arangodb.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	database, err := client.Database(ctx, cfg.Arango.DbName)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	// Define the collections you want to ensure exist
	collections := []string{"users", "manga", "chapter"}

	// Iterate over the collections, checking if they exist and creating them if they don't
	for _, colName := range collections {
		exists, err := database.CollectionExists(context.Background(), colName)
		if err != nil {
			log.Fatalf("Failed to check if collection '%s' exists: %v", colName, err)
		}

		if !exists {
			_, err := database.CreateCollection(context.Background(), colName, nil)
			if err != nil {
				log.Fatalf("Failed to create collection '%s': %v", colName, err)
			} else {
				fmt.Printf("Collection '%s' created successfully\n", colName)
			}
		} else {
			fmt.Printf("Collection '%s' already exists\n", colName)
		}
	}
	return database
}

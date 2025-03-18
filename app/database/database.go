package database

import (
	"context"
	"fmt"
	"time"

	"github.com/theotruvelot/books/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Database struct {
	cfg    *config.Config
	client *mongo.Client
}

func NewDatabase(config *config.Config) (*Database, error) {
	uri := fmt.Sprintf("mongodb://%s:%s",
		config.DatabaseConfig.Host,
		config.DatabaseConfig.Port)

	fmt.Printf("Connecting to MongoDB: %s\n", uri)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	for i := 0; i < 5; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := client.Ping(pingCtx, readpref.Primary())
		pingCancel()

		if err == nil {
			fmt.Println("Successfully connected to MongoDB")
			break
		}

		if i == 4 {
			return nil, fmt.Errorf("failed to ping MongoDB after 5 attempts: %w", err)
		}

		fmt.Printf("Failed to ping MongoDB (attempt %d/5): %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return &Database{
		cfg:    config,
		client: client,
	}, nil
}

func (d *Database) GetClient() *mongo.Client {
	return d.client
}

func (d *Database) Close(ctx context.Context) error {
	if d.client != nil {
		return d.client.Disconnect(ctx)
	}
	return nil
}

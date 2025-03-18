package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/theotruvelot/books/api"
	"github.com/theotruvelot/books/books"
	"github.com/theotruvelot/books/config"
	"github.com/theotruvelot/books/database"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	fmt.Printf("Using database: %s\n", config.DatabaseConfig.Database)
	bookRepo := books.NewBookRepository(db.GetClient(), config.DatabaseConfig.Database)

	server := api.NewServer(bookRepo)

	fmt.Printf("Starting server on port %s\n", config.ServerConfig.Port)
	if err := server.Start(":" + config.ServerConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

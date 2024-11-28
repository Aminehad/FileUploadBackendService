// Description: The main entry point for the upload service application.
package main

import (
	"context"
	"log"
	"os"
	"upload-service/internal/app/upload-service/api"
	"upload-service/internal/app/upload-service/repository"
	"upload-service/pkg/app"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// TODO: move this to a configuration file .yml or .env
	// data source name for the postgres database
	PostgresDSN = "postgres://user_file_upload:pass@postgres:5432/file_upload?sslmode=disable"
)

func main() {
	// Connect to the database
	db, err := gorm.Open(postgres.Open(PostgresDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "upload-service: ", log.LstdFlags|log.Lshortfile)

	// Create repository
	repository := repository.New(db)
	if err := repository.AutoMigrate(); err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// Initialize API
	apiService := api.New(api.Configuration{
		Log:        logger,
		Repository: repository,
	})

	// Create and start the app
	app := app.New(apiService)
	defer app.Stop()

	logger.Println("starting upload service")
	if err := app.Start(context.Background()); err != nil {
		logger.Fatalf("failed to start application: %s", err)
	}
}

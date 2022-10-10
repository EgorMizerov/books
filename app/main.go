package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/egormizerov/books/app/config"
	"github.com/egormizerov/books/app/database"
	"github.com/egormizerov/books/app/database/client"
	"github.com/egormizerov/books/app/handlers"
	"github.com/egormizerov/books/app/services"
	"github.com/egormizerov/books/pkg/log"
	"github.com/egormizerov/books/pkg/process"
	"github.com/egormizerov/books/pkg/wrappers"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file")
	}
}

func main() {
	appConfig := config.NewAppConfig()
	logger := log.NewLogrusLogger(appConfig.LoggerEnableJson, appConfig.LoggerLogLevel)
	databaseConnection, err := database.ConnectToDatabase(&wrappers.SimpleSqlxWrapper{}, database.ConnectConfig{
		User:     appConfig.DatabaseUser,
		Password: appConfig.DatabasePassword,
		Host:     appConfig.DatabaseHost,
		Port:     appConfig.DatabasePort,
		Database: appConfig.DatabaseDatabase,
	})
	if err != nil {
		logger.
			WithError(err).
			Fatal("failed to get database connection")
	}

	databaseClient := client.NewDatabaseClient(databaseConnection)
	service := services.NewService(databaseClient, &wrappers.SimpleUUIDWrapper{})
	handler := handlers.NewHandler(logger, service, validator.New())

	go func() {
		if err = http.ListenAndServe(":8080", handler); err != nil {
			logger.
				WithError(err).
				Fatal("Server unexpectedly stopped")
		}
	}()
	logger.Info("Server has started!")
	process.WaitForTermination()
}

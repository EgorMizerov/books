package main

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/egormizerov/books/app/config"
	"github.com/egormizerov/books/app/database"
	"github.com/egormizerov/books/app/database/client"
	"github.com/egormizerov/books/app/handlers"
	"github.com/egormizerov/books/app/services"
	"github.com/egormizerov/books/pkg/log"
	"github.com/egormizerov/books/pkg/process"
	"github.com/egormizerov/books/pkg/server"
	"github.com/egormizerov/books/pkg/wrappers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file")
	}

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
	serverHost := fmt.Sprintf("%s:%s", appConfig.ServerHost, appConfig.ServerPort)
	httpServer := server.NewServer(serverHost, handler)
	go func() {
		if err = httpServer.Listen(); err != nil {
			logger.
				WithError(err).
				Fatal("server unexpectedly stopped")
		}
	}()
	logger.Info("server has started")

	process.WaitForTermination()
	if err = httpServer.Shutdown(context.Background()); err != nil {
		logger.
			WithError(err).
			Fatal("failed to shutdown server")
	}
}

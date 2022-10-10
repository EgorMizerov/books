package config

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/egormizerov/books/pkg/env"
)

const (
	configKeyPrefix = "books"

	configKeyLoggerLogLevel   = configKey("LoggerLogLevel")
	configKeyLoggerEnableJson = configKey("LoggerEnableJson")

	configKeyDatabaseUser     = configKey("DatabaseUser")
	configKeyDatabasePassword = configKey("DatabasePassword")
	configKeyDatabaseHost     = configKey("DatabaseHost")
	configKeyDatabasePort     = configKey("DatabasePort")
	configKeyDatabaseDatabase = configKey("DatabaseDatabase")

	configKeyServerPort = configKey("ServerPort")
	configKeyServerHost = configKey("ServerHost")
)

type configKey string

func (self configKey) String() string {
	return fmt.Sprintf("%s_%s", configKeyPrefix, string(self))
}

type AppConfig struct {
	LoggerLogLevel   logrus.Level
	LoggerEnableJson bool

	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseDatabase string

	ServerPort string
	ServerHost string
}

func NewAppConfig() AppConfig {
	return AppConfig{
		LoggerLogLevel:   logrus.Level(env.GetInt(configKeyLoggerLogLevel.String(), int(logrus.InfoLevel))),
		LoggerEnableJson: env.GetBool(configKeyLoggerEnableJson.String(), false),

		DatabaseUser:     env.GetString(configKeyDatabaseUser.String(), "postgres"),
		DatabasePassword: env.GetString(configKeyDatabasePassword.String(), "postgres"),
		DatabaseHost:     env.GetString(configKeyDatabaseHost.String(), "localhost"),
		DatabasePort:     env.GetString(configKeyDatabasePort.String(), "5432"),
		DatabaseDatabase: env.GetString(configKeyDatabaseDatabase.String(), "postgres"),

		ServerPort: env.GetString(configKeyServerPort.String(), "8080"),
		ServerHost: env.GetString(configKeyServerHost.String(), "localhost"),
	}
}

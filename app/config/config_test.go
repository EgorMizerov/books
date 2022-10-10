package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type AppConfigTests struct {
	suite.Suite
}

func TestAppConfig(t *testing.T) {
	suite.Run(t, new(AppConfigTests))
}

func (self *AppConfigTests) TestNewAppConfig() {
	loggerLogLevel := logrus.DebugLevel
	loggerEnableJson := true
	databaseUser := "test_user"
	databasePassword := "test_password"
	databaseHost := "test_host"
	databasePort := "test_port"
	databaseDatabase := "test_database"
	self.NoError(os.Setenv(configKeyLoggerLogLevel.String(), strconv.Itoa(int(loggerLogLevel))))
	self.NoError(os.Setenv(configKeyLoggerEnableJson.String(), strconv.FormatBool(loggerEnableJson)))
	self.NoError(os.Setenv(configKeyDatabaseUser.String(), databaseUser))
	self.NoError(os.Setenv(configKeyDatabasePassword.String(), databasePassword))
	self.NoError(os.Setenv(configKeyDatabaseHost.String(), databaseHost))
	self.NoError(os.Setenv(configKeyDatabasePort.String(), databasePort))
	self.NoError(os.Setenv(configKeyDatabaseDatabase.String(), databaseDatabase))

	result := NewAppConfig()

	self.Equal(AppConfig{
		LoggerLogLevel:   loggerLogLevel,
		LoggerEnableJson: loggerEnableJson,
		DatabaseUser:     databaseUser,
		DatabasePassword: databasePassword,
		DatabaseHost:     databaseHost,
		DatabasePort:     databasePort,
		DatabaseDatabase: databaseDatabase,
	}, result)
}

func (self *AppConfigTests) TestConfigKeyToString() {
	configKeyValue := "test_value"
	configKey := configKey(configKeyValue)

	actualResult := configKey.String()

	expectedResult := fmt.Sprintf("%s_%s", configKeyPrefix, configKeyValue)
	self.Equal(expectedResult, actualResult)
}

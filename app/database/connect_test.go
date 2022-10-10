package database

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/egormizerov/books/pkg/wrappers"
	wrappersmocks "github.com/egormizerov/books/pkg/wrappers/mocks"
)

type ConnectToDatabaseTests struct {
	suite.Suite
	sqlxWrapper wrappers.SqlxWrapper
	sqlxMock    *wrappersmocks.SqlxWrapper
}

func TestConnectToDatabase(t *testing.T) {
	suite.Run(t, new(ConnectToDatabaseTests))
}

func (self *ConnectToDatabaseTests) SetupTest() {
	self.sqlxMock = wrappersmocks.NewSqlxWrapper(self.T())
	self.sqlxWrapper = self.sqlxMock
}

func (self *ConnectToDatabaseTests) TestConnectToDatabaseErrorIfNotConnect() {
	connectionConfig := ConnectConfig{}
	connectionError := errors.New("test_error")
	self.sqlxMock.On("Connect", "pgx", getConnectionUrl(connectionConfig)).
		Return(nil, connectionError)

	database, err := ConnectToDatabase(self.sqlxWrapper, connectionConfig)

	self.EqualError(connectionError, err.Error())
	self.Nil(database)
}

func (self *ConnectToDatabaseTests) TestConnectToDatabase() {
	connectionConfig := ConnectConfig{}
	self.sqlxMock.On("Connect", "pgx", getConnectionUrl(connectionConfig)).
		Return(&sqlx.DB{}, nil)

	result, err := ConnectToDatabase(self.sqlxWrapper, connectionConfig)

	self.Nil(err)
	self.Equal(&sqlx.DB{}, result)
}

func (self *ConnectToDatabaseTests) TestGetConnectionUrl() {
	connectionConfig := ConnectConfig{
		User:     "test_user",
		Password: "test_password",
		Host:     "test_host",
		Port:     "test_port",
		Database: "test_database",
	}

	actualResult := getConnectionUrl(connectionConfig)

	expectedResult := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		connectionConfig.User, connectionConfig.Password, connectionConfig.Host, connectionConfig.Port, connectionConfig.Database)
	self.Equal(expectedResult, actualResult)
}

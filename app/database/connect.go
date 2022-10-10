package database

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/egormizerov/books/pkg/wrappers"
)

type ConnectConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func ConnectToDatabase(sqlxWrapper wrappers.SqlxWrapper, config ConnectConfig) (*sqlx.DB, error) {
	url := getConnectionUrl(config)
	conn, err := sqlxWrapper.Connect("pgx", url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getConnectionUrl(config ConnectConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, config.Database)
}

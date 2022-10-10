package wrappers

import "github.com/jmoiron/sqlx"

type SqlxWrapper interface {
	Connect(driverName string, dataSourceName string) (*sqlx.DB, error)
}

type SimpleSqlxWrapper struct{}

func (self *SimpleSqlxWrapper) Connect(driverName string, dataSourceName string) (*sqlx.DB, error) {
	return sqlx.Connect(driverName, dataSourceName)
}

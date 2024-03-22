package database

import (
	"fmt"
	"github.com/justjack1521/mevconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrFailedConnectToPostgres = func(err error) error {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
)

func NewPostgresConnection() (*gorm.DB, error) {
	config, err := mevconn.NewPostgresConfig()
	if err != nil {
		return nil, ErrFailedConnectToPostgres(err)
	}
	db, err := gorm.Open(postgres.Open(config.Source()), &gorm.Config{})
	if err != nil {
		return nil, ErrFailedConnectToPostgres(err)
	}
	return db, nil
}

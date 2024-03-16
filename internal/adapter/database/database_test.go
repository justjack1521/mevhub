package database_test

import (
	"github.com/jinzhu/configor"
	"gorm.io/gorm"
	"mevhub/internal/config"
	"os"
)

const conf string = "/src/mevhub/internal/config/config.test.json"

func NewDatabaseConnection() *gorm.DB {
	var path = os.Getenv("GOPATH")
	var configuration config.Application
	if err := configor.Load(&configuration, path+"/"+conf); err != nil {
		panic(err)
	}
	conn, err := configuration.Database.NewPostgresConnection()
	if err != nil {
		panic(err)
	}
	return conn
}

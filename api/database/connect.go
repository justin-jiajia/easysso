package database

import (
	"log"

	"github.com/justin-jiajia/easysso/api/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getdatabase() gorm.Dialector {
	if config.Config.DBType == "sqlite" {
		return sqlite.Open(config.Config.DBPath)
	} else if config.Config.DBType == "mysql" {
		return mysql.Open(config.Config.DBPath)
	} else {
		log.Panicf("Unknown database type: %v", config.Config.DBType)
		return nil
	}
}

func InitDB() {
	var err error
	DB, err = gorm.Open(getdatabase(), &gorm.Config{})
	if err != nil {
		log.Panicf("Cannot open the database! Error: %v", err.Error())
	}
}

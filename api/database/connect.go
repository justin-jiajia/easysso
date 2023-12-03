package database

import (
	"log"

	"github.com/justin-jiajia/easysso/api/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.Config.DBSavePath), &gorm.Config{})
	if err != nil {
		log.Panicf("Cannot open the database! Error: %v", err.Error())
	}
}

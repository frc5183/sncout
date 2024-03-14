package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	models2 "sncout/models"
)

var db *gorm.DB

var models = []interface{}{
	&models2.Team{},
	&models2.Match{},
	&models2.Robot{},
}

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("sncout.db"))
	if err != nil {
		panic("failed to connect database")
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			panic("failed to migrate model")
		}
	}
}

func GetDB() *gorm.DB {
	return db
}

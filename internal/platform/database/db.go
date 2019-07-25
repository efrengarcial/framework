package database

import (
	"github.com/jinzhu/gorm"
)

func Initialize(dbDriver string, dbURI string) *gorm.DB {

	// Get database details from environment variables
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		panic("failed to connect database")
	}
	db.LogMode(true)

	return db
}

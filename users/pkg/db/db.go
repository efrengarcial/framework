package db

import (
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/jinzhu/gorm"
)

func Initialize(dbDriver string, dbURI string) *gorm.DB {

	// Get database details from environment variables
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		panic("failed to connect database")
	}
	db.LogMode(true)

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&model.User{}, &model.Authority{}, &model.Privilege{})

	return db
}


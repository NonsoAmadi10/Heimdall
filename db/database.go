package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("metrics.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return db

}

package db

import (
	"log"

	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("metrics.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.AutoMigrate(&utils.ConnectionMetrics{})

	return db

}

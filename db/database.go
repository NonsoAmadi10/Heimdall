package db

import (
	"fmt"

	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("metrics.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.AutoMigrate(&utils.ConnectionMetrics{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil

}

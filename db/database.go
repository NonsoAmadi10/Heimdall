package db

import (
	"fmt"
	"strings"
	"sync"

	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	initErr  error
	once     sync.Once
)

func DB() (*gorm.DB, error) {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open("metrics.db"), &gorm.Config{})
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		if err := db.AutoMigrate(&utils.ConnectionMetrics{}, &utils.Alert{}); err != nil && !strings.Contains(err.Error(), "already exists") {
			initErr = fmt.Errorf("failed to migrate database: %w", err)
			return
		}

		instance = db
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

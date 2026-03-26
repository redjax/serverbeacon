package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
	// Extract dir path from db path
	dir := filepath.Dir(path)
	// Ensure parent dir exists
	if err := os.Mkdir(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed creating db directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on", path)

	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gdb, nil
}

// Close ensures the connection to the database is closed.
// Useful during app shutdown.
func Close(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()

	if err != nil {
		return err
	}

	return sqlDB.Close()
}

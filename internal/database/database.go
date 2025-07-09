package database

import (
	"github.com/cliffdoyle/go-auth-app/internal/model"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect initializes the database connection and runs migrations.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	slog.Info("Running database migrations...")
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}
	slog.Info("Database migration completed.")

	return db, nil
}

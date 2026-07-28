package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-chatbot/internal/models"
)

type DB struct {
	GormDB *gorm.DB
}

func Open(dsn string) (*DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := gormDB.AutoMigrate(&models.Conversation{}); err != nil {
		return nil, err
	}

	return &DB{GormDB: gormDB}, nil
}

func (db *DB) Close() error {
	sqlDB, err := db.GormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

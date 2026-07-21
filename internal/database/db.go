package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"go-chatbot/internal/models"
)

type DB struct {
	GormDB *gorm.DB
}

func Open(path string) (*DB, error) {
	dsn := path + "?_pragma=journal_mode(WAL)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=synchronous(NORMAL)"

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := gormDB.AutoMigrate(&models.Conversation{}, &models.ConversationMessage{}); err != nil {
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

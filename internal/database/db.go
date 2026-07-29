package database

import (
	"embed"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-chatbot/internal/models"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	GormDB *gorm.DB
}

func Open(dsn string) (*DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	if err := gormDB.AutoMigrate(&models.Conversation{}); err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
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

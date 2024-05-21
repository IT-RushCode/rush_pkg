package database

import (
	"fmt"

	"github.com/IT-RushCode/rush_pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PSQL_CONNECT(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Pass, cfg.Name, cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("database is not connected!")
	}

	return db
}

func PostgresClose(db *gorm.DB) error {
	dbDB, err := db.DB()
	if err != nil {
		return err
	}

	return dbDB.Close()
}

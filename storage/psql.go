package storage

import (
	"fmt"

	"github.com/IT-RushCode/rush_pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PSQL_CONNECT создает подключение к PostgreSQL с возможностью передать кастомный GORM-конфиг
// и произвольные хуки для до- или пост-конфигурации *gorm.DB.
func PSQL_CONNECT(
	cfg *config.DatabaseConfig,
	gormCfg *gorm.Config,
	hooks ...func(*gorm.DB) error,
) *gorm.DB {
	if gormCfg == nil {
		gormCfg = &gorm.Config{}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.HOST, cfg.USER, cfg.PASS, cfg.NAME, cfg.PORT,
	)

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		panic("database is not connected: " + err.Error())
	}

	for _, hook := range hooks {
		if hook == nil {
			continue
		}
		if err := hook(db); err != nil {
			panic("failed to execute hook: " + err.Error())
		}
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

package storage

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Storage struct {
	PSQL  *gorm.DB
	REDIS *redis.Client
}

// Флаги для указания инициализации баз данных
const (
	InitPSQL  = 1 << iota // 1
	InitREDIS             // 2
)

var (
	storageInstance *Storage
	storageOnce     sync.Once
)

// Инициализация баз данных в соответствии с переданными флагами
func DB_CONNECT(
	cfg *config.Config,
	initFlags int,
	gormCfg *gorm.Config,
	redisModifiers []func(*redis.Options),
	hooks ...func(*gorm.DB) error,
) *Storage {
	storageOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var psqlDB *gorm.DB
		var redisClient *redis.Client

		if initFlags&InitPSQL != 0 {
			poolHook := buildPoolHook(&cfg.DB)
			hooksToApply := make([]func(*gorm.DB) error, 0, len(hooks)+1)
			if poolHook != nil {
				hooksToApply = append(hooksToApply, poolHook)
			}
			hooksToApply = append(hooksToApply, hooks...)

			psqlDB = PSQL_CONNECT(&cfg.DB, gormCfg, hooksToApply...)
		}

		if initFlags&InitREDIS != 0 {
			redisClient = REDIS_CONNECT(ctx, &cfg.REDIS, redisModifiers...)
		}

		storageInstance = &Storage{
			PSQL:  psqlDB,
			REDIS: redisClient,
		}
	})

	return storageInstance
}

func GetDatabaseInstance(
	cfg *config.Config,
	initFlags int,
	gormCfg *gorm.Config,
	redisModifiers []func(*redis.Options),
	hooks ...func(*gorm.DB) error,
) *Storage {
	return DB_CONNECT(cfg, initFlags, gormCfg, redisModifiers, hooks...)
}

func buildPoolHook(cfg *config.DatabaseConfig) func(*gorm.DB) error {
	if cfg == nil {
		return nil
	}

	return func(db *gorm.DB) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		if cfg.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		}
		if cfg.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		}
		if cfg.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		}
		if cfg.ConnMaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
		}

		return nil
	}
}

func CloseDatabases(storage *Storage) {
	if storage.PSQL != nil {
		if db, err := storage.PSQL.DB(); err == nil {
			if err := db.Close(); err != nil {
				log.Printf("Error closing PSQL connection: %v", err)
			}
		} else {
			log.Printf("Error retrieving underlying SQL connection: %v", err)
		}
	}

	if storage.REDIS != nil {
		REDIS_CLOSE(storage.REDIS)
	}
}

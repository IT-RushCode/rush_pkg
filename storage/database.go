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
func DB_CONNECT(cfg *config.Config, initFlags int) *Storage {
	storageOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var psqlDB *gorm.DB
		var redisClient *redis.Client

		if initFlags&InitPSQL != 0 {
			psqlDB = PSQL_CONNECT(&cfg.DB)
		}

		if initFlags&InitREDIS != 0 {
			redisClient = REDIS_CONNECT(ctx, &cfg.REDIS)
		}

		storageInstance = &Storage{
			PSQL:  psqlDB,
			REDIS: redisClient,
		}
	})

	return storageInstance
}

func GetDatabaseInstance(cfg *config.Config, initFlags int) *Storage {
	return DB_CONNECT(cfg, initFlags)
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

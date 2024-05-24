package database

import (
	"sync"

	"github.com/IT-RushCode/rush_pkg/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Storage struct {
	PSQL  *gorm.DB
	REDIS *redis.Client
}

var (
	storageInstance *Storage
	storageOnce     sync.Once
)

func DB_CONNECT(cfg *config.Config) *Storage {
	storageOnce.Do(func() {
		storageInstance = &Storage{
			PSQL:  PSQL_CONNECT(&cfg.DB),
			REDIS: REDIS_CONNECT(&cfg.REDIS),
		}
	})

	return storageInstance
}

func GetDatabaseInstance(cfg *config.Config) *Storage {
	return DB_CONNECT(cfg)
}

package database

import (
	"sync"

	"github.com/IT-RushCode/rush_pkg/config"
	"gorm.io/gorm"
)

type Storage struct {
	MYSQL *gorm.DB
	PSQL  *gorm.DB
	REDIS interface{}
}

var (
	storageInstance *Storage
	storageOnce     sync.Once
)

func DB_CONNECT(cfg *config.Config) *Storage {
	storageOnce.Do(func() {
		storageInstance = &Storage{
			MYSQL: MYSQL_CONNECT(&cfg.DB.MYSQL),
			PSQL:  PSQL_CONNECT(&cfg.DB.PSQL),
			REDIS: REDIS_CONNECT(&cfg.REDIS),
		}
	})

	return storageInstance
}

func GetDatabaseInstance(cfg *config.Config) *Storage {
	return DB_CONNECT(cfg)
}

package database

import (
	"sync"
	"time"

	"gitlab.arvand.tj/conveyor/arvand_pkg/config"
	"gorm.io/gorm"
)

type Storage struct {
	MSSQL          *gorm.DB
	MYSQL          *gorm.DB
	PSQL           *gorm.DB
	MSSQL_ABS      *gorm.DB
	MSSQL_CONVEYOR *gorm.DB
	REDIS          interface{}
}

var (
	storageInstance *Storage
	storageOnce     sync.Once
)

func DB_CONNECT(cfg *config.Config) *Storage {
	storageOnce.Do(func() {
		storageInstance = &Storage{
			MSSQL:          MSSQL_CONNECT(&cfg.DB.MSSQL),
			MYSQL:          MYSQL_CONNECT(&cfg.DB.MYSQL),
			PSQL:           PSQL_CONNECT(&cfg.DB.PSQL),
			MSSQL_ABS:      MSSQL_CONNECT(&cfg.DB.MSSQL_ABS),
			MSSQL_CONVEYOR: MSSQL_CONNECT(&cfg.DB.MSSQL_CONVEYOR),
			REDIS:          REDIS_CONNECT(&cfg.REDIS),
		}
	})

	return storageInstance
}

func GetDatabaseInstance(cfg *config.Config) *Storage {
	if cfg.DB.MSSQL.Host != "localhost" {
		time.Sleep(7 * time.Second)
	}

	return DB_CONNECT(cfg)
}

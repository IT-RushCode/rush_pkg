package storage

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Storage struct {
	PSQL  *gorm.DB
	REDIS *redis.Client
	MONGO *mongo.Client
}

// Флаги для указания инициализации баз данных
const (
	InitPSQL  = 1 << iota // 1
	InitREDIS             // 2
	InitMONGO             // 4
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
		var mongoClient *mongo.Client

		if initFlags&InitPSQL != 0 {
			psqlDB = PSQL_CONNECT(&cfg.DB)
		}

		if initFlags&InitREDIS != 0 {
			redisClient = REDIS_CONNECT(ctx, &cfg.REDIS)
		}

		if initFlags&InitMONGO != 0 {
			mongoClient = MONGO_DB_CONNECT(ctx, &cfg.MONGODB)
		}

		storageInstance = &Storage{
			PSQL:  psqlDB,
			REDIS: redisClient,
			MONGO: mongoClient,
		}
	})

	return storageInstance
}

func GetDatabaseInstance(cfg *config.Config, initFlags int) *Storage {
	return DB_CONNECT(cfg, initFlags)
}

func CloseDatabases(storage *Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	if storage.MONGO != nil {
		if err := MONGO_DB_CLOSE(ctx, storage.MONGO); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}
}

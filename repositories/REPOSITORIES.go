package repositories

import (
	rpBase "github.com/IT-RushCode/rush_pkg/repositories/base"
	rpChat "github.com/IT-RushCode/rush_pkg/repositories/chat"
	rpNtf "github.com/IT-RushCode/rush_pkg/repositories/notification"
	rpYKassa "github.com/IT-RushCode/rush_pkg/repositories/yookassa"
	"github.com/IT-RushCode/rush_pkg/storage"

	"github.com/redis/go-redis/v9"
)

// Флаги для определения, какие репозитории инициализировать
type RepoFlags struct {
	InitYKassaRepo       bool // Инициализировать ли YKassa репозиторий
	InitNotificationRepo bool // Инициализировать ли Notification репозиторий
	InitChatRepo         bool // Инициализировать ли Chat репозиторий
	InitCacheRepo        bool // Инициализировать ли кэш-репозиторий Redis
	InitMongoRepo        bool // Инициализировать ли MongoDB репозиторий
}

// Все репозитории
type Repositories struct {
	YooKassaSetting rpYKassa.YooKassaSettingRepository
	Notification    rpNtf.NotificationRepository
	Chat            rpChat.ChatRepository

	Redis *redis.Client

	Mongo rpBase.MongoBaseRepository
}

// Инициализация всех репозиториев с учетом переданных флагов
func NewRepositories(
	db *storage.Storage,
	flags RepoFlags,
	mongoDB string,
) *Repositories {
	DB := db.PSQL

	repos := &Repositories{}

	// Инициализация репозиториев для YooKassaSetting
	if flags.InitYKassaRepo {
		repos.YooKassaSetting = rpYKassa.NewYooKassaSettingRepository(DB)
	}

	// Инициализация репозиториев для Notification
	if flags.InitNotificationRepo {
		repos.Notification = rpNtf.NewNotificationRepository(DB)
	}

	// Инициализация репозиториев для Chat
	if flags.InitChatRepo {
		repos.Chat = rpChat.NewChatRepository(DB)
	}

	// Инициализация кэш-репозитория Redis
	if flags.InitCacheRepo {
		repos.Redis = redis.NewClient(db.REDIS.Options())
	}

	// Инициализация MongoDB репозитория
	if flags.InitMongoRepo {
		MONGO := db.MONGO.Database(mongoDB)
		repos.Mongo = rpBase.NewMongoBaseRepository(MONGO)
	}

	return repos
}

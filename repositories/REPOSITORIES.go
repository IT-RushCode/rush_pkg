package repositories

import (
	rpAppVersion "github.com/IT-RushCode/rush_pkg/repositories/app-version"
	rpChat "github.com/IT-RushCode/rush_pkg/repositories/chat"
	rpNtf "github.com/IT-RushCode/rush_pkg/repositories/notification"
	rpPolicy "github.com/IT-RushCode/rush_pkg/repositories/policy"
	rpYKassa "github.com/IT-RushCode/rush_pkg/repositories/yookassa"
	"github.com/IT-RushCode/rush_pkg/storage"

	"github.com/redis/go-redis/v9"
)

// Флаги для определения, какие репозитории инициализировать
type RepoFlags struct {
	InitYKassaRepo       bool // Инициализировать ли YKassa репозиторий
	InitPolicyRepo       bool // Инициализировать ли Policy репозиторий
	InitNotificationRepo bool // Инициализировать ли Notification репозиторий
	InitChatRepo         bool // Инициализировать ли Chat репозиторий
	InitCacheRepo        bool // Инициализировать ли кэш-репозиторий Redis
	InitAppVersion       bool // Инициализировать ли AppVersion репозиторий
}

// Все репозитории
type Repositories struct {
	YooKassaSetting rpYKassa.YooKassaSettingRepository
	Notification    rpNtf.NotificationRepository
	Chat            rpChat.ChatRepository
	Policy          rpPolicy.PolicyRepository
	AppVersion      rpAppVersion.AppVersionRepository

	Redis *redis.Client
}

// Инициализация всех репозиториев с учетом переданных флагов
func NewRepositories(
	db *storage.Storage,
	flags RepoFlags,
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

	// Инициализация репозиториев для Policy
	if flags.InitPolicyRepo {
		repos.Policy = rpPolicy.NewPolicyRepository(DB)
	}

	// Инициализация репозиториев для AppVersion
	if flags.InitAppVersion {
		repos.AppVersion = rpAppVersion.NewAppVersionRepository(DB)
	}

	// Инициализация кэш-репозитория Redis
	if flags.InitCacheRepo {
		repos.Redis = db.REDIS
	}

	return repos
}

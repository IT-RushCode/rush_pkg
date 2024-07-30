package repositories

import (
	"github.com/IT-RushCode/rush_pkg/database"
	rpAuth "github.com/IT-RushCode/rush_pkg/repositories/auth"
	rpBase "github.com/IT-RushCode/rush_pkg/repositories/base"
	rpYKassa "github.com/IT-RushCode/rush_pkg/repositories/yookassa"
	"github.com/redis/go-redis/v9"
)

// Флаги для определения, какие репозитории инициализировать
type RepoFlags struct {
	InitAuthRepo   bool // Инициализировать ли репозитории для авторизации
	InitYKassaRepo bool // Инициализировать ли YKassa репозиторий
	InitCacheRepo  bool // Инициализировать ли кэш-репозиторий Redis
	InitMongoRepo  bool // Инициализировать ли MongoDB репозиторий
}

// Все репозитории
type Repositories struct {
	User           rpAuth.UserRepository
	Role           rpAuth.RoleRepository
	UserRole       rpAuth.UserRoleRepository
	Permission     rpAuth.PermissionRepository
	RolePermission rpAuth.RolePermissionRepository

	YooKassaSetting rpYKassa.YooKassaSettingRepository

	Redis *redis.Client

	Mongo rpBase.MongoBaseRepository
}

// Инициализация всех репозиториев с учетом переданных флагов
func NewRepositories(db *database.Storage, flags RepoFlags, mongoDB string) *Repositories {
	DB := db.PSQL

	repos := &Repositories{}

	// Инициализация репозиториев для авторизации
	if flags.InitAuthRepo {
		repos.User = rpAuth.NewUserRepository(DB)
		repos.Role = rpAuth.NewRoleRepository(DB)
		repos.UserRole = rpAuth.NewUserRoleRepository(DB)
		repos.Permission = rpAuth.NewPermissionRepository(DB)
		repos.RolePermission = rpAuth.NewRolePermissionRepository(DB)
	}

	// Инициализация репозиториев для авторизации
	if flags.InitYKassaRepo {
		repos.YooKassaSetting = rpYKassa.NewYooKassaSettingRepository(DB)
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

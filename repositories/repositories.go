package repositories

import (
	"github.com/IT-RushCode/rush_pkg/database"
	rpAuth "github.com/IT-RushCode/rush_pkg/repositories/auth"
	"github.com/redis/go-redis/v9"
)

// Все репозитории
type Repositories struct {
	User       rpAuth.UserRepository
	Role       rpAuth.RoleRepository
	Permission rpAuth.PermissionRepository

	Redis *redis.Client
}

// Инициализация всех репозиториев
func NewRepositories(db *database.Storage) *Repositories {
	DB := db.PSQL

	return &Repositories{
		// rush_pkg repos
		User:       rpAuth.NewUserRepository(DB),
		Role:       rpAuth.NewRoleRepository(DB),
		Permission: rpAuth.NewPermissionRepository(DB),

		// Cache
		Redis: redis.NewClient(db.REDIS.Options()),
	}
}

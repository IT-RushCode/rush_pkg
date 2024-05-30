package repositories

import (
	"github.com/IT-RushCode/rush_pkg/database"
	rpAuth "github.com/IT-RushCode/rush_pkg/repositories/auth"
	rpReview "github.com/IT-RushCode/rush_pkg/repositories/review"
)

// Все репозитории
type Repositories struct {
	User       rpAuth.UserRepository
	Role       rpAuth.RoleRepository
	Permission rpAuth.PermissionRepository
	Review     rpReview.ReviewRepository
}

// Инициализация всех репозиториев
func NewRepositories(db *database.Storage) *Repositories {
	DB := db.PSQL

	return &Repositories{
		// rush_pkg repos
		User:       rpAuth.NewUserRepository(DB),
		Role:       rpAuth.NewRoleRepository(DB),
		Permission: rpAuth.NewPermissionRepository(DB),
		Review:     rpReview.NewReviewRepository(DB),
	}
}

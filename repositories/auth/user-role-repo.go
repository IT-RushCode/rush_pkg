package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/IT-RushCode/rush_pkg/models/auth"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

type UserRoleRepository interface {
	rp.BaseRepository

	BindUserRoles(ctx context.Context, userID uint, newRoles []uint) error
}

type userRoleRepository struct {
	db *gorm.DB
	rp.BaseRepository
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

func (repo *userRoleRepository) BindUserRoles(ctx context.Context, userID uint, newRoles []uint) error {
	// Копируем роли в новый слайс
	newRoleIDs := make([]uint, len(newRoles))
	copy(newRoleIDs, newRoles)

	// Удаляем старые роли, кроме тех, которые переданы в input
	if err := repo.db.WithContext(ctx).Where("user_id = ? AND role_id NOT IN (?)", userID, newRoleIDs).Delete(&auth.UserRole{}).Error; err != nil {
		return err
	}

	// Добавляем новые роли и игнорируем ошибки отсутствующих ролей
	for _, roleID := range newRoleIDs {
		userRole := auth.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := repo.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).FirstOrCreate(&userRole).Error; err != nil {
			// Проверяем, что ошибка связана с отсутствующей ролью
			if strings.Contains(err.Error(), "foreign key constraint fails") || strings.Contains(err.Error(), "violates foreign key constraint") {
				fmt.Printf("Skipping role ID %d due to foreign key constraint failure\n", roleID)
				continue
			}
			return err
		}
	}

	return nil
}

package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/IT-RushCode/rush_pkg/models/auth"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

type RolePermissionRepository interface {
	rp.BaseRepository

	BindRolePermissions(ctx context.Context, roleID uint, newPermissions []uint) error
}

type rolePermissionRepository struct {
	db *gorm.DB
	rp.BaseRepository
}

func NewRolePermissionRepository(db *gorm.DB) RolePermissionRepository {
	return &rolePermissionRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

func (repo *rolePermissionRepository) BindRolePermissions(ctx context.Context, roleID uint, newPermissions []uint) error {
	// Копируем привилегии в новый слайс
	newPermissionIDs := make([]uint, len(newPermissions))
	copy(newPermissionIDs, newPermissions)

	// Удаляем старые привилегии, кроме тех, которые переданы в input
	if err := repo.db.WithContext(ctx).Where("role_id = ? AND permission_id NOT IN (?)", roleID, newPermissionIDs).Delete(&auth.RolePermission{}).Error; err != nil {
		return err
	}

	// Добавляем новые привилегии и игнорируем ошибки отсутствующих привилегий
	for _, permissionID := range newPermissionIDs {
		rolePermission := auth.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
		if err := repo.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).FirstOrCreate(&rolePermission).Error; err != nil {
			// Проверяем, что ошибка связана с отсутствующей привилегией
			if strings.Contains(err.Error(), "foreign key constraint fails") || strings.Contains(err.Error(), "violates foreign key constraint") {
				fmt.Printf("Skipping permission ID %d due to foreign key constraint failure\n", permissionID)
				continue
			}
			return err
		}
	}

	return nil
}

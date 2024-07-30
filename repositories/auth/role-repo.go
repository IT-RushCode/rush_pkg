package repositories

import (
	"context"

	"github.com/IT-RushCode/rush_pkg/models/auth"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

type RoleRepository interface {
	rp.BaseRepository

	FindByIDWithPermissions(ctx context.Context, id uint) (*auth.Role, error)
}

type roleRepository struct {
	db *gorm.DB
	rp.BaseRepository
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

func (repo *roleRepository) FindByIDWithPermissions(ctx context.Context, id uint) (*auth.Role, error) {
	role := &auth.Role{}
	if err := repo.BaseRepository.FindByID(ctx, id, role); err != nil {
		return nil, err
	}

	var permissions auth.Permissions
	if err := repo.db.WithContext(ctx).
		Table(`"Permissions"`).
		Select(`"Permissions".*`).
		Joins(`JOIN "RolePermissions" ON "RolePermissions".permission_id = "Permissions".id`).
		Where(`"RolePermissions".role_id = ?`, role.ID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	role.Permissions = permissions

	return role, nil
}

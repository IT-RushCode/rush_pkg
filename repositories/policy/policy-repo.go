package repositories

import (
	"context"

	"github.com/IT-RushCode/rush_pkg/models"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

type PolicyRepository interface {
	rp.BaseRepository
	FindByKey(ctx context.Context, key string) (*models.Policy, error)
	UpdateText(ctx context.Context, policyType, text string) error
}

type policyRepository struct {
	rp.BaseRepository
	db *gorm.DB
}

func NewPolicyRepository(db *gorm.DB) PolicyRepository {
	return &policyRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

// FindByKey implements PolicyRepository.
func (repo *policyRepository) FindByKey(ctx context.Context, key string) (*models.Policy, error) {
	policy := &models.Policy{Key: key}
	if err := repo.db.WithContext(ctx).
		Find(&policy).
		Error; err != nil {
		return nil, err
	}
	return policy, nil
}

// Обновление текста политики
func (repo *policyRepository) UpdateText(ctx context.Context, policyType, text string) error {
	err := repo.db.WithContext(ctx).
		Model(&models.Policy{}).
		Where("key = ?", policyType).
		Update("text", text).Error
	if err != nil {
		return err
	}
	return nil
}

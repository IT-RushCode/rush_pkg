package repositories

import (
	"context"
	"errors"

	rpDTO "github.com/IT-RushCode/rush_pkg/dto/auth"
	rpModels "github.com/IT-RushCode/rush_pkg/models/auth"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"
	"github.com/IT-RushCode/rush_pkg/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	rp.BaseRepository
	FindByUsernameAndPassword(ctx context.Context, data rpDTO.AuthWithLoginPasswordRequestDTO) (*rpModels.User, error)
	FindByEmailAndPassword(ctx context.Context, data rpDTO.AuthWithEmailPasswordRequestDTO) (*rpModels.User, error)
	FindByPhone(ctx context.Context, data rpDTO.AuthWithPhoneRequestDTO) (*rpModels.User, error)
}

type userRepository struct {
	db *gorm.DB
	rp.BaseRepository
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

var (
	ErrUsernameOrPassword = errors.New("неверный логин или пароль")
	ErrEmailOrPassword    = errors.New("неверный email или пароль")
)

// Аутентификация пользователя по логин-паролю.
func (repo *userRepository) FindByUsernameAndPassword(
	ctx context.Context,
	data rpDTO.AuthWithLoginPasswordRequestDTO,
) (*rpModels.User, error) {
	var user *rpModels.User

	// Поиск пользователя по username
	if err := repo.db.WithContext(ctx).
		Where("user_name = ?", data.Username).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUsernameOrPassword
		}
		return nil, err
	}

	// Сравниваем хеши пароля из запроса с паролем из базы
	if !utils.ComparePassword(user.HashPassword, data.Password) {
		return nil, ErrUsernameOrPassword
	}

	return user, nil
}

// Аутентификация пользователя по логин-паролю.
func (repo *userRepository) FindByEmailAndPassword(
	ctx context.Context,
	data rpDTO.AuthWithEmailPasswordRequestDTO,
) (*rpModels.User, error) {
	var user *rpModels.User

	// Поиск пользователя по email
	if err := repo.db.WithContext(ctx).
		Where("email = ?", data.Email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmailOrPassword
		}
		return nil, err
	}

	// Сравниваем хеши пароля из запроса с паролем из базы
	if !utils.ComparePassword(user.HashPassword, data.Password) {
		return nil, ErrEmailOrPassword
	}

	return user, nil
}

// Аутентификация пользователя по логин-паролю.
func (repo *userRepository) FindByPhone(
	ctx context.Context,
	data rpDTO.AuthWithPhoneRequestDTO,
) (*rpModels.User, error) {
	var user *rpModels.User

	// Поиск пользователя по username
	if err := repo.db.WithContext(ctx).
		Where("phone_number = ?", data.PhoneNumber).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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
	ChangePassword(ctx context.Context, userID uint, dto rpDTO.ChangePasswordRequestDTO) error
	ResetPassword(ctx context.Context, userID uint, newPassword string) (string, error)

	// Cron task
	CheckUserBirthDays(ctx context.Context) error
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
	ErrOldPassword        = errors.New("неверный старый пароль")
	ErrConfirmNewPassword = errors.New("новый пароль и подтверждение не совпадают")
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
	if !utils.ComparePasswordWithSalt(user.HashPassword, data.Password, user.Salt) {
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
	if !utils.ComparePasswordWithSalt(user.HashPassword, data.Password, user.Salt) {
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

// ChangePassword изменяет пароль пользователя.
func (repo *userRepository) ChangePassword(ctx context.Context, userID uint, dto rpDTO.ChangePasswordRequestDTO) error {
	var user rpModels.User

	// Поиск пользователя по ID
	if err := repo.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		return err
	}

	// Проверка старого пароля с использованием соли
	if !utils.ComparePasswordWithSalt(user.HashPassword, dto.OldPassword, user.Salt) {
		return ErrOldPassword
	}

	// Проверка совпадения нового пароля и его подтверждения
	if dto.NewPassword != dto.ConfirmPassword {
		return ErrConfirmNewPassword
	}

	// Хеширование нового пароля с солью
	salt := utils.GenerateSalt()
	hashedPassword, err := utils.HashPasswordWithSalt(dto.NewPassword, salt)
	if err != nil {
		return err
	}

	// Обновление пароля и соли пользователя
	user.HashPassword = hashedPassword
	user.Salt = salt
	user.ChangePasswordWhenLogin = new(bool)
	if err := repo.db.WithContext(ctx).Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword сбрасывает пароль пользователя.
func (repo *userRepository) ResetPassword(ctx context.Context, userID uint, newPassword string) (string, error) {
	var user rpModels.User

	// Поиск пользователя по ID
	if err := repo.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utils.ErrRecordNotFound
		}
		return "", err
	}

	// Хеширование нового пароля с солью
	salt := utils.GenerateSalt()
	hashedPassword, err := utils.HashPasswordWithSalt(newPassword, salt)
	if err != nil {
		return "", err
	}

	// Обновление пароля и соли пользователя
	user.HashPassword = hashedPassword
	user.Salt = salt
	changePassword := true
	user.ChangePasswordWhenLogin = &changePassword
	if err := repo.db.WithContext(ctx).Save(&user).Error; err != nil {
		return "", err
	}

	return user.Email, nil
}

func (repo *userRepository) CheckUserBirthDays(ctx context.Context) error {
	var users rpModels.Users
	today := time.Now().Format("2006-01-02")

	if err := repo.db.WithContext(ctx).
		Where("TO_CHAR(birth_date, 'MM-DD') = TO_CHAR(CAST(? AS DATE), 'MM-DD')", today).
		Find(&users).Error; err != nil {
		return err
	}

	if len(users) > 0 {
		for _, user := range users {
			message := fmt.Sprintf("%s, с Днем рождения!", user.FirstName)

			log.Println(message) // TODO: тут реализовать логику отправки уведомления через notification
		}
	}

	log.Println("Некого поздравлять")

	return nil
}

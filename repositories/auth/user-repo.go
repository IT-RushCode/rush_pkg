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
	FindByUsernameAndPassword(ctx context.Context, data rpDTO.AuthWithLoginPasswordRequestDTO) (*rpModels.User, []string, error)
	FindByEmailAndPassword(ctx context.Context, data rpDTO.AuthWithEmailPasswordRequestDTO) (*rpModels.User, []string, error)
	FindByPhone(ctx context.Context, data rpDTO.AuthWithPhoneRequestDTO) (*rpModels.User, error)
	FindByIDWithRoles(ctx context.Context, id uint) (*rpModels.User, []string, error)
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
) (*rpModels.User, []string, error) {
	var user *rpModels.User

	// Поиск пользователя по username
	if err := repo.db.WithContext(ctx).
		Where("user_name = ?", data.Username).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUsernameOrPassword
		}
		return nil, nil, err
	}

	// Сравниваем хеши пароля из запроса с паролем из базы
	if !utils.ComparePasswordWithSalt(user.HashPassword, data.Password, user.Salt) {
		return nil, nil, ErrUsernameOrPassword
	}

	permissions, err := repo.getRolesAndPermissions(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, permissions, nil
}

// Аутентификация пользователя по логин-паролю.
func (repo *userRepository) FindByEmailAndPassword(
	ctx context.Context,
	data rpDTO.AuthWithEmailPasswordRequestDTO,
) (*rpModels.User, []string, error) {
	var user *rpModels.User

	// Поиск пользователя по email
	if err := repo.db.WithContext(ctx).
		Where("email = ?", data.Email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrEmailOrPassword
		}
		return nil, nil, err
	}

	// Сравниваем хеши пароля из запроса с паролем из базы
	if !utils.ComparePasswordWithSalt(user.HashPassword, data.Password, user.Salt) {
		return nil, nil, ErrEmailOrPassword
	}

	permissions, err := repo.getRolesAndPermissions(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, permissions, nil
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
func (repo *userRepository) ChangePassword(
	ctx context.Context,
	userID uint,
	dto rpDTO.ChangePasswordRequestDTO,
) error {
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
func (repo *userRepository) ResetPassword(
	ctx context.Context,
	userID uint,
	newPassword string,
) (string, error) {
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

func (repo *userRepository) FindByIDWithRoles(ctx context.Context, id uint) (*rpModels.User, []string, error) {
	user := &rpModels.User{}
	if err := repo.BaseRepository.FindByID(ctx, id, user); err != nil {
		return nil, nil, err
	}

	permissions, err := repo.getRolesAndPermissions(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, permissions, nil
}

// Получение всех ролей и привилегий
func (repo *userRepository) getRolesAndPermissions(ctx context.Context, user *rpModels.User) ([]string, error) {
	var roles []rpModels.Role
	if err := repo.db.WithContext(ctx).
		Table(`"Roles"`).
		Select(`"Roles".*`).
		Joins(`JOIN "UserRoles" ON "UserRoles".role_id = "Roles".id`).
		Where(`"UserRoles".user_id = ?`, user.ID).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	user.Roles = roles

	// Маппинг для удаления дубликатов привилегий
	permissionsMap := make(map[string]struct{})

	// Получение всех привилегий для ролей
	for _, role := range roles {
		var permissions []string
		if err := repo.db.WithContext(ctx).
			Table(`"Permissions"`).
			Select(`"Permissions".name`).
			Joins(`JOIN "RolePermissions" ON "RolePermissions".permission_id = "Permissions".id`).
			Where(`"RolePermissions".role_id = ?`, role.ID).
			Pluck("name", &permissions).Error; err != nil {
			return nil, err
		}

		// Добавление привилегий в маппинг, удаляя дубликаты
		for _, permission := range permissions {
			permissionsMap[permission] = struct{}{}
		}
	}

	// Преобразование маппинга в слайс
	uniquePermissions := make([]string, 0, len(permissionsMap))
	for permission := range permissionsMap {
		uniquePermissions = append(uniquePermissions, permission)
	}

	return uniquePermissions, nil
}

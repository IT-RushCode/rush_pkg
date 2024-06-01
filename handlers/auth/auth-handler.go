package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	rpAuthDTO "github.com/IT-RushCode/rush_pkg/dto/auth"
	rpModels "github.com/IT-RushCode/rush_pkg/models/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

var (
	ErrRefreshToken    = errors.New("неверный токен обновления")
	ErrNotRefreshToken = errors.New("полученный токен не является refresh токеном")
)

type AuthHandler struct {
	repo       *repositories.Repositories
	cfg        *config.Config
	jWTService utils.JWTService
	ttl        uint64
	rttl       uint64
}

func NewAuthHandler(
	repo *repositories.Repositories,
	cfg *config.Config,
) *AuthHandler {
	ttl := time.Duration(uint64(cfg.JWT.JWT_TTL)) * time.Second
	rttl := time.Duration(uint64(cfg.JWT.REFRESH_TTL)) * time.Second
	jwtService := utils.NewJWTService(cfg.JWT.JWT_SECRET, ttl, rttl)
	return &AuthHandler{
		repo:       repo,
		cfg:        cfg,
		jWTService: jwtService,
		ttl:        uint64(cfg.JWT.JWT_TTL),
		rttl:       uint64(cfg.JWT.REFRESH_TTL),
	}
}

// Авторизация пользователя
func (h *AuthHandler) PhoneLogin(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.AuthWithPhoneRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	repoRes, err := h.repo.User.FindByPhone(context.Background(), *input)
	if err != nil {
		return utils.ErrorNotFoundResponse(ctx, err.Error(), nil)
	}

	// Обновляем даты последней активности "LastActivity"
	h.updateLastActivity(context.Background(), repoRes.ID)

	accessToken, refreshToken, err := h.jWTService.GenerateTokens(
		repoRes.ID,
		repoRes.FirstName,
		repoRes.UserName,
		repoRes.IsPersonal,
	)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	userRes := &rpAuthDTO.UserPhoneDataDTO{}
	if err := copier.Copy(&userRes, &repoRes); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &rpAuthDTO.AuthPhoneResponseDTO{
		Token: &rpAuthDTO.TokenResponseDTO{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiredIn:  h.ttl,
			RefreshTokenExpiredIn: h.rttl,
		},
		User: userRes,
	}

	return utils.SuccessResponse(ctx, "success", res)
}

// Авторизация пользователя
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.AuthWithLoginPasswordRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	repoRes, err := h.repo.User.FindByUsernameAndPassword(context.Background(), *input)
	if err != nil {
		return utils.ErrorUnauthorizedResponse(ctx, err.Error(), nil)
	}
	if !repoRes.IsPersonal {
		return utils.ErrorForbiddenResponse(ctx, "нет прав", nil)
	}

	// Обновляем даты последней активности "LastActivity"
	h.updateLastActivity(context.Background(), repoRes.ID)

	accessToken, refreshToken, err := h.jWTService.GenerateTokens(
		repoRes.ID,
		repoRes.FirstName,
		repoRes.UserName,
		repoRes.IsPersonal,
	)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	userRes := &rpAuthDTO.UserResponseDTO{}
	if err := copier.Copy(&userRes, &repoRes); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &rpAuthDTO.AuthResponseDTO{
		Token: &rpAuthDTO.TokenResponseDTO{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiredIn:  h.ttl,
			RefreshTokenExpiredIn: h.rttl,
		},
		User: userRes,
	}

	return utils.SuccessResponse(ctx, "success", res)
}

// Получение данных пользователя по ID из токена
func (h *AuthHandler) Me(ctx *fiber.Ctx) error {
	// Получение UserID из локальных данных контекста
	userID := ctx.Locals("UserID").(uint)

	// Получение данных пользователя по UserID
	data := &rpModels.User{}
	if err := h.repo.User.FindByID(context.Background(), userID, data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	// Обновляем даты последней активности "LastActivity"
	h.updateLastActivity(context.Background(), data.ID)

	// Возврат информации
	userDTO := &rpAuthDTO.UserResponseDTO{}
	if err := copier.Copy(userDTO, data); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "success", userDTO)
}

// Обновление токена
func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.RefreshTokenDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	valid, err := h.jWTService.ValidateToken(input.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)

	}

	if isRefreshToken, ok := valid.IsRefreshToken.(bool); !ok || !isRefreshToken {
		return utils.ErrorResponse(ctx, ErrNotRefreshToken.Error(), nil)
	}
	if valid != nil {
		newAccessToken, newRefreshToken, err := h.jWTService.GenerateTokens(
			valid.UserID,
			valid.Name,
			valid.Login,
			valid.IsPersonal,
		)
		if err != nil {
			return err
		}

		tokenDto := &rpAuthDTO.TokenResponseDTO{
			AccessToken:           newAccessToken,
			RefreshToken:          newRefreshToken,
			AccessTokenExpiredIn:  h.ttl,
			RefreshTokenExpiredIn: h.rttl,
		}

		return utils.SuccessResponse(ctx, "success", tokenDto)
	}

	return utils.ErrorResponse(ctx, ErrRefreshToken.Error(), nil)

}

func (h *AuthHandler) updateLastActivity(ctx context.Context, userID uint) {
	go func() {
		if err := h.repo.User.UpdateField(ctx, userID, "last_activity", time.Now(), &rpModels.User{}); err != nil {
			fmt.Printf("Ошибка обновления LastActivity: %v\n", err)
		}
	}()
}

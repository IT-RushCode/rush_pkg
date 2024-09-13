package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	dto "github.com/IT-RushCode/rush_pkg/dto/sms"
	"github.com/IT-RushCode/rush_pkg/services"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type SmsHandler struct {
	cfg   *config.Config
	srv   *services.Services
	cache *redis.Client
}

func NewSMSHandler(
	cfg *config.Config,
	srv *services.Services,
	cache *redis.Client,
) *SmsHandler {
	return &SmsHandler{
		cfg:   cfg,
		srv:   srv,
		cache: cache,
	}
}

// SendSMS обрабатывает запрос на отправку SMS
func (h *SmsHandler) SendSMS(ctx *fiber.Ctx) error {
	var req dto.SMSRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	req.Messages[0].Phone = strings.TrimSpace(req.Messages[0].Phone)

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	res, err := h.srv.Sms.SendSMS(h.cfg, req)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	if res.Message.Status == "error" {
		return utils.ErrorInternalServerErrorResponse(ctx, utils.ErrInternal.Error(), nil)
	}

	// Если OTP код успешно отправлен то сохраняем его в кеше для дальнейшей верификации
	if req.IsOTP && res.Message.Data[0].Status == "sent" {
		err = h.cache.Set(ctx.Context(), res.Phone, res.OTPCode, 5*time.Minute).Err() // Установка времени истечения в 5 минут
		if err != nil {
			return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при сохранении OTP кода в кеш: "+err.Error(), nil)
		}
		return utils.SuccessResponse(ctx, utils.Success, fmt.Sprintf("Код подтверждения отправлен на номер %s", res.Phone))
	}

	if len(req.Messages) == 1 {
		return utils.SuccessResponse(ctx, utils.Success, fmt.Sprintf("Сообщение отправлено на номер %s", res.Phone))
	}

	return utils.SuccessResponse(ctx, "Сообщения отправлены", nil)
}

// VerifySMSCode подтверждает SMS код который был отправлен на номер телефона и удаляет из кеша
func (h *SmsHandler) VerifySMSCode(ctx *fiber.Ctx) error {
	var req dto.VerifyRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка парсинга тела запроса: "+err.Error(), nil)
	}
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка валидации данных запроса: "+err.Error(), nil)
	}

	// Проверка OTP кода в Redis
	otp, err := h.cache.Get(ctx.Context(), req.PhoneNumber).Result()
	if err == redis.Nil {
		return utils.ErrorBadRequestResponse(ctx, "Неверный или истекший код подтверждения", nil)
	} else if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении OTP кода из кеша: "+err.Error(), nil)
	}

	if otp != req.OTPCode {
		return utils.ErrorBadRequestResponse(ctx, "Неверный код подтверждения", nil)
	}

	// Удаляем OTP код из Redis
	err = h.cache.Del(ctx.Context(), req.PhoneNumber).Err()
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при удалении OTP кода из кеша: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Код успешно подтвержден", nil)
}

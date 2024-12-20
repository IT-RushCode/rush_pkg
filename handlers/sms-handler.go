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
	// Проверка флага отправки SMS
	if !h.cfg.SMS.ACTIVE_SEND {
		return utils.SuccessResponse(ctx, "Пропуск отправки SMS кода", nil)
	}

	// Разбор входящего запроса
	var req dto.SMSRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}
	req.Messages[0].Phone = strings.TrimSpace(req.Messages[0].Phone)

	// Получение списка игнорируемых номеров
	ignoringNumbers := strings.Split(h.cfg.SMS.IGNORING_NUMBERS, ",")

	// Проверка игнорируемых номеров
	if h.shouldIgnore(req.Messages[0].Phone, req.IsOTP, ignoringNumbers) {
		return utils.SuccessResponse(ctx, fmt.Sprintf("Пропуск отправки для номера %s", req.Messages[0].Phone), nil)
	}

	// Валидация запроса
	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	// Отправка SMS
	res, err := h.srv.Sms.SendSMS(h.cfg, req)
	if err != nil {
		return err
	}

	// Проверка статуса отправки
	if res.Message.Status == "error" {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка отправки SMS на стороне сервера", nil)
	}

	// Сохранение OTP в Redis (если это OTP-код)
	if req.IsOTP && res.Message.Data[0].Status == "sent" {
		if err := h.cache.Set(ctx.Context(), res.Phone, res.OTPCode, 5*time.Minute).Err(); err != nil {
			return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при сохранении OTP кода в кеш", err)
		}
		return utils.SuccessResponse(ctx, fmt.Sprintf("Код подтверждения отправлен на номер %s", res.Phone), nil)
	}

	// Обработка одиночного сообщения
	if len(req.Messages) == 1 {
		return utils.SuccessResponse(ctx, fmt.Sprintf("Сообщение отправлено на номер %s", res.Phone), nil)
	}

	// Ответ для множественных сообщений
	return utils.SuccessResponse(ctx, "Сообщения отправлены", nil)
}

// shouldIgnore проверяет, нужно ли игнорировать отправку SMS для указанного номера
func (h *SmsHandler) shouldIgnore(phone string, isOTP bool, ignoringNumbers []string) bool {
	for _, ignoringPhone := range ignoringNumbers {
		if phone == ignoringPhone && isOTP {
			return true
		}
	}
	return false
}

// VerifySMSCode подтверждает SMS код который был отправлен на номер телефона и удаляет из кеша
func (h *SmsHandler) VerifySMSCode(ctx *fiber.Ctx) error {
	var req dto.VerifyRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	if err := utils.ValidateStruct(req); err != nil {
		return err
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

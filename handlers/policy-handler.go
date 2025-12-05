package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type PolicyHandler struct {
	repo *repositories.Repositories
}

func NewPolicyHandler(repo *repositories.Repositories) *PolicyHandler {
	return &PolicyHandler{repo: repo}
}

var (
	ErrorTypePolicy = errors.New("необходимо указать один из типов политики: privacy, agreement")
)

type policyTextDTO struct {
	Text string `json:"text"`
}

// UpdatePolicyText обновляет текст политики
// @Summary Обновление текста политики
// @Description Обновляет текст политики privacy или agreement
// @Tags Политики
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param policyType path string true "Тип политики" Enums(privacy, agreement)
// @Param data body policyTextDTO true "Новый текст политики"
// @Success 200 {object} utils.Response "Политика обновлена"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /policy/{policyType} [patch]
func (h *PolicyHandler) UpdatePolicyText(ctx *fiber.Ctx) error {
	policyType := ctx.Params("policyType")
	if policyType != "privacy" && policyType != "agreement" {
		return utils.ErrorBadRequestResponse(ctx, ErrorTypePolicy.Error(), nil)
	}

	// Получение нового текста политики из тела запроса
	body := policyTextDTO{}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка парсинга тела запроса: "+err.Error(), nil)
	}

	// Вызов репозитория для обновления текста
	if err := h.repo.Policy.UpdateText(ctx.Context(), policyType, body.Text); err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка обновления текста политики: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, utils.Success, nil)
}

// GetPolicy возвращает JSON-представление политики
// @Summary Получение политики
// @Description Возвращает текст политики privacy или agreement в JSON
// @Tags Политики
// @Produce json
// @Security BearerAuth
// @Param policyType path string true "Тип политики" Enums(privacy, agreement)
// @Success 200 {object} utils.Response "Текст политики"
// @Failure 400 {object} utils.Response "Неверный тип политики"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /policy/{policyType} [get]
func (h *PolicyHandler) GetPolicy(ctx *fiber.Ctx) error {
	policyType := ctx.Params("policyType")
	if policyType != "privacy" && policyType != "agreement" {
		return utils.ErrorBadRequestResponse(ctx, ErrorTypePolicy.Error(), nil)
	}

	policy, err := h.repo.Policy.FindByKey(ctx.Context(), policyType)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(ctx, utils.Success, policy)
}

// GetPolicyHTML возвращает HTML версию политики (публичная страница)
// @Summary Публичная политика в HTML
// @Description Возвращает HTML текст политики для web-страницы
// @Tags Политики
// @Produce html
// @Param policyType path string true "Тип политики" Enums(privacy, agreement)
// @Success 200 {string} string "HTML контент"
// @Failure 400 {object} utils.Response "Неверный тип политики"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /public/policy/{policyType} [get]
func (h *PolicyHandler) GetPolicyHTML(ctx *fiber.Ctx) error {
	policyType := ctx.Params("policyType")
	if policyType != "privacy" && policyType != "agreement" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid policy type")
	}

	policy, err := h.repo.Policy.FindByKey(ctx.Context(), policyType)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(utils.ErrClientInternal)
	}

	// Рендеринг HTML-шаблона
	return ctx.Render("policy", fiber.Map{
		"Title": policy.Title,
		"Text":  convertToHTML(policy.Text),
	})
}

func convertToHTML(text string) string {
	// Пример: разбиваем текст на абзацы и оборачиваем в <p>
	paragraphs := strings.Split(text, "\n")
	var html string
	for _, p := range paragraphs {
		html += fmt.Sprintf("<p>%s</p>", p)
	}
	return html
}

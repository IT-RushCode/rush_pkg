package utils

import (
	"github.com/IT-RushCode/rush_pkg/dto"
	"github.com/gofiber/fiber/v2"
)

func GetAllQueries(ctx *fiber.Ctx) (*dto.GetAllRequest, error) {
	var req dto.GetAllRequest

	// Парсинг основных параметров
	if err := ctx.QueryParser(&req); err != nil {
		return nil, ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	// Установка значений по умолчанию
	if req.Limit == 0 {
		req.Limit = 20
	}

	if req.Offset == 0 {
		req.Offset = 1
	}

	// Извлечение фильтров вручную
	filters := make(map[string]string)
	queryParams := ctx.Queries()
	for key, value := range queryParams {
		if value != "" {
			if len(key) > 8 && key[:8] == "filters[" && key[len(key)-1] == ']' {
				filterKey := key[8 : len(key)-1]
				filters[filterKey] = value
			}
		}
	}
	req.Filters = filters

	return &req, nil
}

package utils

import "strings"

// IsRouteMatch проверяет, соответствует ли путь маршруту, включая поддержку параметров.
func IsRouteMatch(path, route string) bool {
	// Удаляем завершающие слеши для маршрута и пути
	path = strings.TrimRight(path, "/")
	route = strings.TrimRight(route, "/")

	// Если оба пустые после удаления слешей, значит это корневой маршрут
	if path == "" && route == "" {
		return true
	}

	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	// Если количество частей не совпадает, маршруты не соответствуют
	if len(routeParts) != len(pathParts) {
		return false
	}

	// Сравниваем каждую часть маршрута
	for i := range routeParts {
		// Если это динамическая часть, пропускаем проверку
		if strings.HasPrefix(routeParts[i], ":") {
			continue
		}
		// Если части не совпадают, возвращаем false
		if routeParts[i] != pathParts[i] {
			return false
		}
	}

	// Если все части совпали, возвращаем true
	return true
}

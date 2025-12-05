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

	// // Если количество частей не совпадает, маршруты не соответствуют
	// if len(routeParts) != len(pathParts) {
	// 	return false
	// }

	// Сравниваем каждую часть маршрута
	for i := range routeParts {
		// Если маршрут содержит wildcard "*", он матчит всё, что идёт дальше.
		if routeParts[i] == "*" {
			return true
		}

		// Если путь закончился раньше, значит не совпадает
		if i >= len(pathParts) {
			return false
		}

		// Если это динамическая часть, пропускаем проверку
		if strings.HasPrefix(routeParts[i], ":") {
			continue
		}
		// Если части не совпадают, возвращаем false
		if routeParts[i] != pathParts[i] {
			return false
		}
	}

	// Если маршрут короче или длиннее пути, он не соответствует,
	// потому что до этого момента все части совпали.
	return len(routeParts) == len(pathParts)
}

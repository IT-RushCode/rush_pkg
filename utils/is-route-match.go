package utils

import "strings"

// IsRouteMatch проверяет, соответствует ли путь маршруту, включая поддержку параметров.
func IsRouteMatch(path, route string) bool {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) != len(pathParts) {
		return false
	}

	for i := range routeParts {
		if routeParts[i] != pathParts[i] && !strings.HasPrefix(routeParts[i], ":") {
			return false
		}
	}
	return true
}

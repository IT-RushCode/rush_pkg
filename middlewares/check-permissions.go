package middlewares

// import (
// 	"strings"

// 	"github.com/IT-RushCode/rush_pkg/models/auth"
// )

// // Получение требуемых привилегий для маршрута
// func (m *AuthMiddleware) getRequiredPermissions(path string) []string {
// 	for route, perms := range m.required {
// 		if strings.HasPrefix(path, route) {
// 			return perms
// 		}
// 	}
// 	return nil
// }

// // Проверка наличия требуемых привилегий у пользователя
// func (m *AuthMiddleware) hasRequiredPermissions(userPerms auth.Permissions, requiredPerms []string) bool {
// 	permsMap := make(map[string]struct{})
// 	for _, perm := range userPerms {
// 		permsMap[perm.Name] = struct{}{}
// 	}

// 	for _, required := range requiredPerms {
// 		if _, exists := permsMap[required]; !exists {
// 			return false
// 		}
// 	}

// 	return true
// }

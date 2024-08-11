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

// // AuthorizationMiddleware проверяет права пользователя
// func AuthorizationMiddleware(requiredPermission string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Предполагается, что вы уже получили userID из токена или сессии
// 		userID := c.Locals("userID")
// 		if userID == nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
// 		}

// 		var user *rpModels.User
// 		if err := db.Preload("Permissions").Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "пользователь не найден"})
// 		}

// 		// Проверка прав пользователя
// 		hasPermission := false
// 		for _, perm := range user.Permissions {
// 			if perm.Name == requiredPermission {
// 				hasPermission = true
// 				break
// 			}
// 		}

// 		// Проверка прав через роли
// 		if !hasPermission {
// 			for _, role := range user.Roles {
// 				for _, perm := range role.Permissions {
// 					if perm.Name == requiredPermission {
// 						hasPermission = true
// 						break
// 					}
// 				}
// 			}
// 		}

// 		if !hasPermission {
// 			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Доступ запрещен"})
// 		}

// 		return c.Next()
// 	}
// }

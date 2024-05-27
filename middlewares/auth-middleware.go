package middlewares

// import (
// 	"encoding/json"
// 	"errors"
// 	"log"

// 	"time"

// 	"github.com/IT-RushCode/rush_pkg/config"
// 	"github.com/IT-RushCode/rush_pkg/utils"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/redis/go-redis/v9"
// )

// var (
// 	ErrAuthHeader = errors.New("missing authorization token").Error()
// 	// ErrAuthToken  = errors.New("token is invalid").Error()
// )

// type AuthData struct {
// 	UserID int64  `json:"userId"`
// 	IP     string `json:"ip"`
// }

// type AuthMiddleware struct {
// 	authClient *cl.AuthClient
// 	jwtTTL     time.Duration
// 	r          *redis.Client
// }

// func NewAuthMiddleware(cfg *config.Config, r *redis.Client) *AuthMiddleware {
// 	authClient, err := cl.NewAuthClient(cfg)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	return &AuthMiddleware{
// 		authClient: authClient,
// 		jwtTTL:     time.Duration(cfg.JWT.JWT_TTL) * time.Second,
// 		r:          r,
// 	}
// }

// func (m *AuthMiddleware) VerifyToken(ctx *fiber.Ctx) error {
// 	// Список маршрутов, которые не требуют проверки токена
// 	noAuthRoutes := []string{
// 		"/",
// 		"/api/v1/auth/login",
// 		"/api/v1/auth/refresh-token",
// 	}
// 	for _, route := range noAuthRoutes {
// 		if ctx.Path() == route {
// 			return ctx.Next()
// 		}
// 	}

// 	// Проверка наличия токена в header
// 	authHeader := ctx.Get("Authorization")
// 	if authHeader == "" {
// 		log.Println(ErrAuthHeader)
// 		return utils.ErrorResponse(ctx, ErrAuthHeader, nil, fiber.StatusUnauthorized)
// 	}

// 	// Проверка токена в редисе
// 	if exists, err := m.r.Exists(ctx.Context(), authHeader).Result(); err != nil {
// 		log.Println(err)
// 	} else if exists == 1 {
// 		// Ключ существует в Redis
// 		return ctx.Next()
// 	}

// 	// Отправляем токен в auth_service на проверку
// 	userId, err := m.authClient.ValidateToken(ctx.Context(), jwt.VerifyToken(authHeader))
// 	if err != nil {
// 		code, msg := utils.HandleGRPCError(err)
// 		return utils.ErrorResponse(
// 			ctx, msg, nil, code,
// 		)
// 	}

// 	clientIP := utils.GetClientIP(ctx)

// 	authData := AuthData{
// 		UserID: userId.GetId(),
// 		IP:     clientIP,
// 	}

// 	authDataJSON, err := json.Marshal(authData)
// 	if err != nil {
// 		log.Println(err)
// 		return utils.ErrorResponse(ctx, "Ошибка при обработке данных", nil, fiber.StatusInternalServerError)
// 	}

// 	// Сохраняем данные в кеше
// 	err = m.r.Set(ctx.Context(), authHeader, authDataJSON, m.jwtTTL).Err()
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	userId.GetId()

// 	return ctx.Next()
// }

package utils

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrorGenAccessToken   = errors.New("не удалось сгенерировать токен доступа")
	ErrorGenRefreshToken  = errors.New("не удалось создать токен обновления")
	ErrorSigningMethod    = errors.New("неверный метод подписи токена")
	ErrorInvalidToken     = errors.New("неверный токен")
	ErrorTokenExpired     = errors.New("токен истёк")
	ErrorTokenNotYetValid = errors.New("токен больше не валидный")
)

type JwtCustomClaim struct {
	UserID         uint   `json:"userId"`
	Name           string `json:"name"`
	Login          string `json:"login"`
	IsPersonal     bool   `json:"isPersonal"`
	IsRefreshToken interface{}
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateTokens(userId uint, name, login string, isPersonal bool) (string, string, error)
	ValidateToken(tokenString string) (*JwtCustomClaim, error)
}

type jwtService struct {
	SecretKey       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

func NewJWTService(secretKey string, accessTokenExp, refreshTokenExp time.Duration) JWTService {
	return &jwtService{
		SecretKey:       secretKey,
		AccessTokenExp:  accessTokenExp,
		RefreshTokenExp: refreshTokenExp,
	}
}

func (service *jwtService) GenerateTokens(userId uint, name, login string, isPersonal bool) (string, string, error) {
	accessTokenExp := time.Now().Add(service.AccessTokenExp)
	refreshTokenExp := time.Now().Add(service.RefreshTokenExp)

	accessTokenClaims := &JwtCustomClaim{
		UserID:         userId,
		Name:           name,
		Login:          login,
		IsPersonal:     isPersonal,
		IsRefreshToken: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
		},
	}

	refreshTokenClaims := &JwtCustomClaim{
		UserID:         userId,
		Name:           name,
		Login:          login,
		IsPersonal:     isPersonal,
		IsRefreshToken: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(service.SecretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(service.SecretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (service *jwtService) ValidateToken(tokenString string) (*JwtCustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return []byte(service.SecretKey), nil
		default:
			return nil, ErrorSigningMethod
		}
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaim)
	if !ok || !token.Valid {
		return nil, ErrorInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrorTokenExpired
	}

	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(time.Now()) {
		return nil, ErrorTokenNotYetValid
	}

	return claims, nil
}

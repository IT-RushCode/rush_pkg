package auth

// Авторизация по телефону
type AuthWithPhoneRequestDTO struct {
	PhoneNumber string `json:"phoneNumber"  validate:"required, length:11"`
}

// Авторизация по логин паролю
type AuthWithLoginPasswordRequestDTO struct {
	Username string `json:"username"  validate:"required, min:3"`
	Password string `json:"password" validate:"required, min:8"`
}

// Ответ авторизации
type AuthResponseDTO struct {
	Token *TokenResponseDTO
	Data  interface{}
}

// Токен
type TokenResponseDTO struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccsessTokenExpiredIn uint64 `json:"accsessTokenExpiredIn"`
	RefreshTokenExpiredIn uint64 `json:"refreshTokenExpiredIn"`
}

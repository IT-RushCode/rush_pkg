package auth

// Авторизация по телефону
type AuthWithPhoneRequestDTO struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,phone,len=12"`
}

type AuthPhoneResponseDTO struct {
	Token *TokenResponseDTO
	User  *UserPhoneDataDTO
}

// Авторизация по логин паролю
type AuthWithLoginPasswordRequestDTO struct {
	Username string `json:"userName" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=8"`
}

// Ответ авторизации
type AuthResponseDTO struct {
	Token *TokenResponseDTO
	User  *UserResponseDTO
}

// Токен
type TokenResponseDTO struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExpiredIn  uint64 `json:"accessTokenExpiredIn"`
	RefreshTokenExpiredIn uint64 `json:"refreshTokenExpiredIn"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refreshToken"`
}

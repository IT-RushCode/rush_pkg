package dto

type AuthWithPhoneRequestDTO struct {
	PhoneNumber string `json:"phoneNumber"  validate:"required, length:11"`
}

type AuthWithLoginPasswordRequestDTO struct {
	Username string `json:"username"  validate:"required, min:3"`
	Password string `json:"password" validate:"required, min:8"`
}

type AuthResponseDTO struct {
	Token *TokenResponseDTO
	Data  interface{}
}

type TokenResponseDTO struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccsessTokenExpiredIn uint64 `json:"accsessTokenExpiredIn"`
	RefreshTokenExpiredIn uint64 `json:"refreshTokenExpiredIn"`
}

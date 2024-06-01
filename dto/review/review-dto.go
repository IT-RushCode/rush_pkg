package dto

import "time"

type ReviewRequestDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required"` // пример
}

type ReviewResponseDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required"` // пример

	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type ReviewsResponseDTO []ReviewResponseDTO

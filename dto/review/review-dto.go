package dto

import "time"

type ReviewRequestDTO struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"userID"`
	PointID uint   `json:"pointID"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
	Status  *bool  `json:"status"`
}

type ReviewResponseDTO struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"userID"`
	PointID uint   `json:"pointID"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
	Status  *bool  `json:"status"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type ReviewsResponseDTO []ReviewResponseDTO

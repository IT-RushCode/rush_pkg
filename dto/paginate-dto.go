package dto

type PaginationDTO struct {
	Data       interface{} `json:"data"`
	Meta       MetaDTO     `json:"meta"`
	TotalCount int64       `json:"totalCount"`
}

type MetaDTO struct {
	Limit  uint `json:"limit"`
	Offset uint `json:"offset"`
}

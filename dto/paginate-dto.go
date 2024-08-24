package dto

type PaginationDTO struct {
	List interface{} `json:"list"`
	Meta MetaDTO     `json:"meta"`
}

type MetaDTO struct {
	Limit      uint  `json:"limit"`
	Offset     uint  `json:"offset"`
	TotalCount int64 `json:"totalCount"`
}

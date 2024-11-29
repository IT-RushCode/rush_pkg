package dto

type GetAllRequest struct {
	Offset  uint              `query:"offset"`
	Limit   uint              `query:"limit"`
	SortBy  string            `query:"sortBy"`
	OrderBy string            `query:"orderBy"`
	Filters map[string]string `query:"f"`
}

type PaginationDTO struct {
	List interface{} `json:"list"`
	Meta MetaDTO     `json:"meta"`
}

type MetaDTO struct {
	Limit      uint  `json:"limit"`
	Offset     uint  `json:"offset"`
	TotalCount int64 `json:"totalCount"`
}

package dto

type GetAllRequest struct {
	Offset  uint              `query:"offset"`
	Limit   uint              `query:"limit"`
	SortBy  string            `query:"sortBy"`
	OrderBy string            `query:"orderBy"`
	Filters map[string]string `query:"f"`
}

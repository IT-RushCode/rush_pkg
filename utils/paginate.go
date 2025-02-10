package utils

import (
	"gorm.io/gorm"
)

// calculatePagination рассчитывает значения limit и offset
func calculatePagination(offset, limit uint) (int, int) {
	// Значения по умолчанию
	if limit == 0 {
		limit = 20
	}

	if offset == 0 {
		offset = 1
	}

	page := int(offset)
	if page <= 0 {
		page = 1
	}

	pageSize := int(limit)

	// Ограничиваем размер страницы
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 20
	default:
	}

	sqlOffset := (page - 1) * pageSize

	return pageSize, sqlOffset
}

// Paginate для GORM-запросов
func Paginate(offset, limit uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize, sqlOffset := calculatePagination(offset, limit)
		return db.Offset(sqlOffset).Limit(pageSize)
	}
}

// SQLPaginate для формирования LIMIT и OFFSET в SQL-запросах
func SQLPaginate(offset, limit uint) (string, []interface{}) {
	pageSize, sqlOffset := calculatePagination(offset, limit)
	sqlClause := "LIMIT ? OFFSET ?"
	args := []interface{}{pageSize, sqlOffset}
	return sqlClause, args
}

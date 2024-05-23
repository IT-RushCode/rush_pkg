package utils

import (
	"gorm.io/gorm"
)

func Paginate(offset, limit uint) func(db *gorm.DB) *gorm.DB {
	if offset == 0 {
		offset = 1
	}

	if limit == 0 {
		limit = 20
	}

	return func(db *gorm.DB) *gorm.DB {
		page := int(offset)
		if page <= 0 {
			page = 1
		}

		pageSize := int(limit)

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 20
		}

		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}

package database

import "gorm.io/gorm"

func NotDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = ?", false)
}

package database

import "gorm.io/gorm"

func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()

	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

package database

import (
	"github.com/daparadoks/go-fuel-api/internal/consumption"
	"github.com/daparadoks/go-fuel-api/internal/member"
	"github.com/jinzhu/gorm"
)

func MigrateDB(db *gorm.DB) error {
	if result := db.AutoMigrate(&member.Member{}); result.Error != nil {
		return result.Error
	}
	if result := db.AutoMigrate(&member.MemberToken{}); result.Error != nil {
		return result.Error
	}

	if result := db.AutoMigrate(&consumption.Consumption{}); result.Error != nil {
		return result.Error
	}

	return nil
}

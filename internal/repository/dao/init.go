package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{}) // 自动迁移数据库
}

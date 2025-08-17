package db

import "gorm.io/gorm"

func InitDB() *gorm.DB {
	Connect()
	Migrate()
	return DB
}

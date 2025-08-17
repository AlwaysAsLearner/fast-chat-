package db

import (
	"log"

	config "github.com/AlwaysAsLearner/fast-chat/backend/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := config.DB.ToDsn()
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect the database :", err)
	}
	DB = database
}

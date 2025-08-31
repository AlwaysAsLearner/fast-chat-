package db

import (
	"log"

	model "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
)

func Migrate() {
	log.Println("Running Migrations")
	err := DB.AutoMigrate(&model.User{}, &model.Chatroom{}, &model.ChatroomMember{}, &model.Message{})
	if err != nil {
		log.Fatal("Migration Failed ", err)
	}
}

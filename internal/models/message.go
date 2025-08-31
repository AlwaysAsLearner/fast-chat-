package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ChatroomID uint      `json:"chatroom_id"`
	SenderID   uint      `json:"sender_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`

	Chatroom Chatroom `gorm:"foreignKey:ChatroomID"`
	Sender   User     `gorm:"foreignKey:SenderID"`
}

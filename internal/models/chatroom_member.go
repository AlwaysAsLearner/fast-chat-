package model

import (
	"time"
	"gorm.io/gorm"
)

type ChatroomMember struct {
	gorm.Model
	UserID     uint      `gorm:"not null;index"`
	ChatroomID uint      `gorm:"not null;index"`
	JoinedAt   time.Time `gorm:"not null;autoCreateTime"`
	IsAdmin    bool      `gorm:"not null;default:false"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID"`
	Chatroom Chatroom `gorm:"foreignKey:ChatroomID"`
}
package model

import (
	"time"
	"gorm.io/gorm"
)

type ChatroomMember struct {
	gorm.Model 
	ID uint `gorm:"primaryKey"`
	UserId uint `gorm:"not null"`
	ChatroomId uint `gorm:"not null"`
	JoinedAt time.Time `gorm:"not null"`
	IsAdmin bool `gorm:"default:false; not null"`
}



package model 

import (
	"gorm.io/gorm"
)

type Chatroom struct {
	gorm.Model
	Name        string     `gorm:"type:varchar(100);not null"`
	Description *string    `gorm:"type:text"`
	IsPrivate   bool       `gorm:"not null;default:false"`
	Members     []ChatroomMember `gorm:"foreignKey:ChatroomID"`
} 


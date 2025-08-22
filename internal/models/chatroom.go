package model 

import (
	"time"

	"gorm.io/gorm"
)

type Chatroom struct {
	gorm.Model 
	ID uint `gorm:"primaryKey"` 
	Name string `gorm:"type:varchar(100); not null"` 
	Description *string 
	IsPrivate bool `gorm:"not null"` 
	CreatedAt time.Time 
	UpdatedAt time.Time 
}
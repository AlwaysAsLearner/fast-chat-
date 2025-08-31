package repositories

import (
	model "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (mr *MessageRepository) Create(message *model.Message) error {
	return mr.DB.Create(message).Error
}

func (mr *MessageRepository) GetByChatroom(chatroomID uint, limit, offset int) ([]model.Message, error) {
	var messages []model.Message
	err := mr.DB.Where("chatroom_id = ?", chatroomID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Preload("Sender").
		Find(&messages).Error

	return messages, err
}

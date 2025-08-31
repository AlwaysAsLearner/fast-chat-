package services

import (
	"time"

	model "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/repositories"
)

type MessageService struct {
	Repo *repositories.MessageRepository
}

func NewMessageService(repo *repositories.MessageRepository) *MessageService {
	return &MessageService{
		Repo: repo,
	}
}

func (ms *MessageService) SendMessage(content string, userID, chatroomID uint) (*model.Message, error) {
	msg := &model.Message{
		ChatroomID: chatroomID,
		SenderID: userID,
		Content: content,
		CreatedAt: time.Now(), 
	}
	err := ms.Repo.Create(msg)
	return msg, err 
}


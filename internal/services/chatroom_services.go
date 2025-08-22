package services

import (
	"errors"
	"fmt"

	model "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/repositories"
	"gorm.io/gorm"
)

type ChatroomService struct {
	repo repositories.ChatroomRepository
}

func NewChatroomService(repo repositories.ChatroomRepository) *ChatroomService {
	return &ChatroomService{repo: repo}
}

func (crs *ChatroomService) CreateChatroom(name, desc string, isPrivate bool, ownerId uint) error {
	_, err := crs.repo.GetByName(name)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a fresh instance instead of using the nil pointer
			newChatroom := model.Chatroom{
				Name:        name,
				Description: &desc,
				IsPrivate:   isPrivate,
			}

			err = crs.repo.Create(newChatroom)
			if err != nil {
				return err
			}

			err = crs.repo.AddUserToChatroom(newChatroom.ID, ownerId, true)
			if err != nil {
				return err
			}
		}
		// Some other DB error
		return fmt.Errorf("failed to check existing chatroom: %w", err)
	}

	// If no error, chatroom already exists
	return fmt.Errorf("chatroom with name %s already exists", name)
}

func (crs *ChatroomService) JoinPublicChatroom(chatId, userId uint) error {
	return crs.repo.AddUserToChatroom(chatId, userId, false)
}

func (crs *ChatroomService) ListAllPublicChatrooms() ([]model.Chatroom, error) {
	return crs.repo.GetPublicChatrooms()
}

func (crs *ChatroomService) LeaveChatrooom(chatId, userId uint) error {
	return crs.repo.RemoveUserFromChatroom(chatId, userId)
}

func (crs *ChatroomService) GetUserJoinedChatrooms(userId uint) ([]model.Chatroom, error) {
	return crs.repo.GetUserChatrooms(userId)
}
package services

import (
	"errors"
	"fmt"
	"time"

	model "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/repositories"
	"gorm.io/gorm"
)

type ChatroomService struct {
	repo repositories.ChatroomRepository
}

type Chatroom struct {
	ID          uint
	Name        string
	Description *string
	IsPrivate   bool
	OwnerID     uint
	MemberCount int
	CreatedAt   time.Time
}

func NewChatroomService(repo repositories.ChatroomRepository) *ChatroomService {
	return &ChatroomService{repo: repo}
}

func (crs *ChatroomService) CreateChatroom(name string, desc *string, isPrivate bool, ownerId uint) (Chatroom, error) {
	_, err := crs.repo.GetByName(name)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a fresh instance instead of using the nil pointer
			newChatroom := model.Chatroom{
				Name:        name,
				Description: desc,
				IsPrivate:   isPrivate,
			}

			err = crs.repo.Create(newChatroom)
			if err != nil {
				return Chatroom{}, err
			}

			chatroom, err := crs.repo.GetByName(name)
			if err != nil {
				return Chatroom{}, err
			}

			err = crs.repo.AddUserToChatroom(chatroom.ID, ownerId, true)
			if err != nil {
				return Chatroom{}, err
			}

			count, err := crs.repo.GetMemberCount(newChatroom.ID)

			chatroomservice := Chatroom{
				ID:          chatroom.ID,
				Name:        chatroom.Name,
				Description: chatroom.Description,
				MemberCount: int(count),
				OwnerID:     ownerId,
				CreatedAt:   newChatroom.CreatedAt,
			}
			return chatroomservice, nil
		}
		// Some other DB error
		return Chatroom{}, fmt.Errorf("failed to check existing chatroom: %w", err)
	}

	// If no error, chatroom already exists
	return Chatroom{}, fmt.Errorf("chatroom with name %s already exists", name)
}

func (crs *ChatroomService) JoinPublicChatroom(chatId, userId uint) error {
	return crs.repo.AddUserToChatroom(chatId, userId, false)
}

func (crs *ChatroomService) toServiceChatroom(r model.Chatroom) Chatroom {
	owner, _ := crs.repo.GetOwner(r.ID)
	memberCount, _ := crs.repo.GetMemberCount(r.ID)
	return Chatroom{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		MemberCount: int(memberCount),
		OwnerID:     owner.ID,
		CreatedAt:   r.CreatedAt,
	}
}

func (crs *ChatroomService) ListAllPublicChatrooms(page, limit int, query string) ([]Chatroom, int64, error) {
	rooms, total, err := crs.repo.GetPublicChatrooms(page, limit, query)
	list := make([]Chatroom, 0, len(rooms))
	for _, r := range rooms {
		list = append(list, crs.toServiceChatroom(r))
	}
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (crs *ChatroomService) LeaveChatrooom(chatId, userId uint) error {
	return crs.repo.RemoveUserFromChatroom(chatId, userId)
}

func (crs *ChatroomService) GetUserJoinedChatrooms(userId uint) ([]Chatroom, error) {
	rooms, err := crs.repo.GetUserChatrooms(userId)
	if err != nil {
		return nil, err
	}
	list := make([]Chatroom, 0, len(rooms))
	for _, r := range rooms {
		list = append(list, crs.toServiceChatroom(r))
	}
	return list, nil
}

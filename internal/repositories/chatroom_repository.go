package repositories

import (
	"fmt"
	"log"
	"strconv"
	"time"

	model "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
	"gorm.io/gorm"
)

type ChatroomRepository struct {
	DB *gorm.DB
}

func NewChatroomRepository(db *gorm.DB) *ChatroomRepository {
	return &ChatroomRepository{
		DB: db,
	}
}

func (c *ChatroomRepository) Create(chatroom model.Chatroom) error {
	return c.DB.Create(&chatroom).Error
}

func (c *ChatroomRepository) GetById(id uint, isPrivate bool) (*model.Chatroom, error) {
	var chatroom model.Chatroom
	err := c.DB.Model(&model.Chatroom{}).Where("id = ? AND is_private=?", id, strconv.FormatBool(isPrivate)).First(&chatroom).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("chatroom with id %d Not found", id)
		}
		log.Printf("Unexpected error while fetching chatroom by %d: %s", id, err)
		return nil, err
	}
	return &chatroom, nil
}

func (c *ChatroomRepository) GetByName(name string) (*model.Chatroom, error) {
	var chatroom model.Chatroom
	err := c.DB.Model(&model.Chatroom{}).Where("name = ?", name).First(&chatroom).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("chatroom with id %s Not found", name)
		}
		log.Printf("Unexpected error while fetching chatroom by %s: %s", name, err)
		return nil, err
	}
	return &chatroom, nil
}

func (c *ChatroomRepository) GetPublicChatrooms() ([]model.Chatroom, error) {
	var chatrooms []model.Chatroom
	err := c.DB.Model(&model.Chatroom{}).Where("is_private = ?", false).Find(&chatrooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no public chatrooms found")
		}
		log.Println("Unexpected error while fetching public chatrooms: ", err)
		return nil, err
	}
	return chatrooms, nil
}

func (c *ChatroomRepository) AddUserToChatroom(chatId, userId uint, isAdmin bool) error {
	var member model.ChatroomMember
	var count int64
	err := c.DB.Model(&model.Chatroom{}).Where("id = ?", chatId).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("chatroom with id %d Not found", chatId)
		}
		return err
	}

	err = c.DB.Model(&model.User{}).Where("id = ? AND is_private=false", userId).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user with id %d Not found", userId)
		}
		return err
	}
	member.ChatroomId = chatId
	member.UserId = userId
	member.IsAdmin = isAdmin
	member.JoinedAt = time.Now().UTC()
	err = c.DB.Create(&member).Error
	return err
}

func (c *ChatroomRepository) RemoveUserFromChatroom(chatId, userId uint) error {
	return c.DB.Where("chatroom_id = ? AND user_id = ?", chatId, userId).Delete(&model.ChatroomMember{}).Error
}

func (c *ChatroomRepository) GetUserChatrooms(userId uint) ([]model.Chatroom, error) {
	var chatrooms []model.Chatroom
	err := c.DB.Model(&model.Chatroom{}).Where("user_id = ? AND is_private=false", userId).Find(&chatrooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no chatrooms found for user with id %d", userId)
		}
		log.Printf("unexpected error while fetching chatrooms for user id %d: %s", userId, err)
		return nil, err
	}
	return chatrooms, nil
}

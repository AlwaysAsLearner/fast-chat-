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

type ChatroomInfo struct {
	ID          uint
	Name        string
	MemberCount int64
	OwnerID     uint
}

func NewChatroomRepository(db *gorm.DB) *ChatroomRepository {
	return &ChatroomRepository{
		DB: db,
	}
}

func (c *ChatroomRepository) Create(chatroom *model.Chatroom) error {
	return c.DB.Model(&model.Chatroom{}).Create(chatroom).Error
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
			return nil, err
		}
		log.Printf("Unexpected error while fetching chatroom by %s: %s", name, err)
		return nil, err
	}
	return &chatroom, nil
}

func (c *ChatroomRepository) GetPublicChatrooms(page, limit int, q string) ([]model.Chatroom, int64, error) {
	var chatrooms []model.Chatroom
	var total int64
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	query := c.DB.Model(&model.Chatroom{}).Where("is_private = ?", false)
	query.Count(&total)
	if q != "" {
		searchTerm := fmt.Sprintf("%%%s%%", q)
		query = query.Where("name ILIKE ?", searchTerm) // ILIKE for case-insensitive search (Postgres)
		query.Count(&total)
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Find(&chatrooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("no public chatrooms found")
		}
		log.Println("Unexpected error while fetching public chatrooms: ", err)
		return nil, 0, err
	}
	return chatrooms, total, nil
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

	err = c.DB.Model(&model.User{}).Where("id = ?", userId).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user with id %d Not found", userId)
		}
		return err
	}
	member.ChatroomID = chatId
	member.UserID = userId
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
	err := c.DB.Model(&model.Chatroom{}).Where("user_id = ? AND is_private = ?", userId, false).Find(&chatrooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no chatrooms found for user with id %d", userId)
		}
		log.Printf("unexpected error while fetching chatrooms for user id %d: %s", userId, err)
		return nil, err
	}
	return chatrooms, nil
}

func (c *ChatroomRepository) GetMemberCount(chatID uint) (int64, error) {
	var count int64
	err := c.DB.Model(&model.ChatroomMember{}).Where("chatroom_id = ?", chatID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *ChatroomRepository) GetOwner(chatID uint) (*model.ChatroomMember, error) {
	var member model.ChatroomMember
	err := c.DB.Model(&model.ChatroomMember{}).Where("chatroom_id = ? AND is_admin = true").First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

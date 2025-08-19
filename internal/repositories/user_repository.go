package repositories

import (
	"fmt"
	"log"

	models "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with email %s Not Found", email)
		}
		log.Println("Unexpected error while fetching user by email:", err)
		return nil, err
	}

	return &user, nil
}

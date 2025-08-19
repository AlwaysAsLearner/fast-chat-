package services

import (
	"fmt"

	model "github.com/AlwaysAsLearner/fast-chat/backend/internal/models"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(email, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		Email:    email,
		Password: string(hash),
		Username: "randomuser",
	}
	err := s.Repo.Create(user)
	return err
}

func (s *AuthService) Authenticate(email, password string) (*model.User, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("password doesn't match for email %s", email)
	}
	return user, nil
}

package service

import (
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/model"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{DB: config.DB}
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

package repository

import (
	"rent-video-game/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	RegisterUser(user *model.Users) (*model.Users, error)
	GetUserByID(userID uuid.UUID) (*model.Users, error)
	GetUserByEmail(email string) (*model.Users, error)
	UpdateUserAmount(userID uuid.UUID, amount int) (*model.Users, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) RegisterUser(user *model.Users) (*model.Users, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	var user model.Users
	if err := r.db.Where("user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		userID, "0001-01-01 00:00:00").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.Users, error) {
	var user model.Users
	if err := r.db.Where("email = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		email, "0001-01-01 00:00:00").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserAmount(userID uuid.UUID, amount int) (*model.Users, error) {
	var user model.Users
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	user.Amount = user.Amount + amount

	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

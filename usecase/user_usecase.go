package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"

	"github.com/google/uuid"
)

type UserUsecase struct {
	userUsecase repository.IUserRepository
}

func NewUserUsecase(userRepo repository.IUserRepository) *UserUsecase {
	return &UserUsecase{userRepo}
}

func (u *UserUsecase) RegisterUser(user *model.Users) (*model.Users, error) {
	var error []string
	if user.Name == "" {
		error = append(error, "Name is required")
	}
	if user.Email == "" {
		error = append(error, "Email is required")
	}
	if !strings.Contains(user.Email, "@") {
		error = append(error, "Email must have @")
	}
	if strings.Contains(user.Email, " ") {
		error = append(error, "Email must not have space")
	}
	if user.Password == "" {
		error = append(error, "Password is required")
	}
	if user.Address == "" {
		error = append(error, "Address is required")
	}
	if user.Amount < 0 {
		error = append(error, "Amount must be greater than 0")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.userUsecase.RegisterUser(user)
}

func (u *UserUsecase) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	return u.userUsecase.GetUserByID(userID)
}

func (u *UserUsecase) GetUserByEmail(email string) (*model.Users, error) {
	return u.userUsecase.GetUserByEmail(email)
}

func (u *UserUsecase) UpdateUserAmount(userID uuid.UUID, amount int) (*model.Users, error) {
	return u.userUsecase.UpdateUserAmount(userID, amount)
}

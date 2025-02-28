package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	RegisterUser(user *model.Users) (*model.Users, error)
	GetUserByID(userID uuid.UUID) (*model.Users, error)
	GetUserByEmail(email string) (*model.Users, error)
	TopupUser(userID uuid.UUID, user *model.Users) (*model.Users, error)
	TransactionUser(userID uuid.UUID, user *model.Users) (*model.Users, error)
}

type UserUsecase struct {
	userRepo repository.IUserRepository
}

func NewUserUsecase(userRepo repository.IUserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) RegisterUser(user *model.Users) (*model.Users, error) {
	var error []string

	if user.Name == "" {
		error = append(error, "name is required")
	}
	if user.Email == "" {
		error = append(error, "email is required")
	}
	if !strings.Contains(user.Email, "@") {
		error = append(error, "email must have @")
	}
	if strings.Contains(user.Email, " ") {
		error = append(error, "email must not have space")
	}
	if user.Password == "" {
		error = append(error, "password is required")
	}
	if user.Address == "" {
		error = append(error, "address is required")
	}
	if user.Amount != 0 {
		error = append(error, "cannot set amount")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.userRepo.RegisterUser(user)
}

func (u *UserUsecase) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	return u.userRepo.GetUserByID(userID)
}

func (u *UserUsecase) GetUserByEmail(email string) (*model.Users, error) {
	return u.userRepo.GetUserByEmail(email)
}

func (u *UserUsecase) TopupUser(userID uuid.UUID, user *model.Users) (*model.Users, error) {
	if user.Amount < 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	return u.userRepo.TopupUser(userID, user)
}

func (u *UserUsecase) TransactionUser(userID uuid.UUID, user *model.Users) (*model.Users, error) {
	if user.Amount < 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	return u.userRepo.TransactionUser(userID, user)
}

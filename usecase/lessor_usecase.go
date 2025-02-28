package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"

	"github.com/google/uuid"
)

type LessorUsecase struct {
	lessorRepo repository.ILessorRepository
}

func NewLessorUsecase(lessorRepo repository.ILessorRepository) *LessorUsecase {
	return &LessorUsecase{lessorRepo: lessorRepo}
}

func (u *LessorUsecase) RegisterLessor(lessor *model.Lessors) (*model.Lessors, error) {
	var error []string

	if lessor.UserID == uuid.Nil {
		error = append(error, "user ID is required")
	}
	if lessor.Name == "" {
		error = append(error, "name is required")
	}
	if lessor.Location == "" {
		error = append(error, "location is required")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.lessorRepo.RegisterLessor(lessor)
}

func (u *LessorUsecase) GetLessorByID(lessorID int) (*model.Lessors, error) {
	return u.lessorRepo.GetLessorByID(lessorID)
}

func (u *LessorUsecase) UpdateLessor(lessorID int, lessor *model.Lessors) (*model.Lessors, error) {
	var error []string
	
	if lessor.UserID != uuid.Nil {
		error = append(error, "user ID cannot be set")
	}
	if lessor.Name == "" {
		error = append(error, "name is required")
	}
	if lessor.Location == "" {
		error = append(error, "location is required")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.lessorRepo.UpdateLessor(lessorID, lessor)
}

func (u *LessorUsecase) DeleteLessor(lessorID int) (*model.Lessors, error) {
	return u.lessorRepo.DeleteLessor(lessorID)
}

func (u *LessorUsecase) GetLessorByUserID(userID uuid.UUID) (*model.Lessors, error) {
	return u.lessorRepo.GetLessorByUserID(userID)
}

func (u *LessorUsecase) GetLessorByProductID(productID int) (*model.Lessors, error) {
	return u.lessorRepo.GetLessorByProductID(productID)
}

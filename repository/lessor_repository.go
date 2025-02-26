package repository

import (
	"errors"
	"rent-video-game/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ILessorRepository interface {
	RegisterLessor(lessor *model.Lessors) (*model.Lessors, error)
	GetLessorByID(lessorID int, userID uuid.UUID) (*model.Lessors, error)
	UpdateLessor(lessorID int, userID uuid.UUID, lessor *model.Lessors) (*model.Lessors, error)
	DeleteLessor(lessorID int, userID uuid.UUID) (*model.Lessors, error)
}

type LessorRepository struct {
	db *gorm.DB
}

func NewLessorRepository(db *gorm.DB) *LessorRepository {
	return &LessorRepository{db}
}

func (r *LessorRepository) RegisterLessor(lessor *model.Lessors) (*model.Lessors, error) {
	var existingLessor model.Lessors
	if err := r.db.Where("user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		lessor.UserID, "0001-01-01 00:00:00").First(&existingLessor).Error; err == nil {
		return nil, errors.New("user already has a registered lessor")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := r.db.Create(lessor).Error; err != nil {
		return nil, err
	}
	return lessor, nil
}

func (r *LessorRepository) GetLessorByID(lessorID int, userID uuid.UUID) (*model.Lessors, error) {
	var lessor model.Lessors
	if err := r.db.Where("lessor_id = ? AND user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		lessorID, userID, "0001-01-01 00:00:00").First(&lessor).Error; err != nil {
		return nil, err
	}
	return &lessor, nil
}

func (r *LessorRepository) UpdateLessor(lessorID int, userID uuid.UUID, lessor *model.Lessors) (*model.Lessors, error) {
	var l model.Lessors
	err := r.db.Where("lessor_id = ? AND user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		lessorID, userID, "0001-01-01 00:00:00").First(&l).Error
	if err != nil {
		return &l, err
	}

	err = r.db.Model(&l).Updates(map[string]interface{}{
		"name":     lessor.Name,
		"location": lessor.Location,
	}).Error
	if err != nil {
		return &l, err
	}

	return &l, nil
}

func (r *LessorRepository) DeleteLessor(lessorID int, userID uuid.UUID) (*model.Lessors, error) {
	var lessor model.Lessors
	if err := r.db.Where("lessor_id = ? AND user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		lessorID, userID, "0001-01-01 00:00:00").First(&lessor).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&lessor).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, err
	}

	return &lessor, nil
}

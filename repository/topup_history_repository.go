package repository

import (
	"rent-video-game/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITopupHistoryRepository interface {
	CreateTopupHistory(topupHistory *model.TopupHistory) (*model.TopupHistory, error)
	GetTopupHistoryByID(topup_history_id int) (*model.TopupHistory, error)
	GetAllTopupHistoryByUser(userID uuid.UUID) (*[]model.TopupHistory, error)
}

type TopupHistoryRepository struct {
	db *gorm.DB
}

func NewTopupHistoryRepository(db *gorm.DB) *TopupHistoryRepository {
	return &TopupHistoryRepository{db}
}

func (r *TopupHistoryRepository) CreateTopupHistory(topupHistory *model.TopupHistory) (*model.TopupHistory, error) {
	if err := r.db.Create(topupHistory).Error; err != nil {
		return nil, err
	}
	return topupHistory, nil
}

func (r *TopupHistoryRepository) GetTopupHistoryByID(topup_history_id int) (*model.TopupHistory, error) {
	var topupHistory model.TopupHistory
	if err := r.db.Where("topup_history_id = ?", topup_history_id).First(&topupHistory).Error; err != nil {
		return nil, err
	}
	return &topupHistory, nil
}

func (r *TopupHistoryRepository) GetAllTopupHistoryByUser(userID uuid.UUID) (*[]model.TopupHistory, error) {
	var topupHistory []model.TopupHistory
	if err := r.db.Where("user_id = ?", userID).Find(&topupHistory).Error; err != nil {
		return nil, err
	}
	return &topupHistory, nil
}

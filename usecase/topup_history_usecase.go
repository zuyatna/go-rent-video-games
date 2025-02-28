package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"

	"github.com/google/uuid"
)

type ITopupHistoryUsecase interface {
	CreateTopupHistory(topupHistory *model.TopupHistory) (*model.TopupHistory, error)
	GetTopupHistoryByID(topup_history_id int) (*model.TopupHistory, error)
	GetAllTopupHistoryByUser(userID uuid.UUID) (*[]model.TopupHistory, error)
}

type TopupHistoryUsecase struct {
	topupHistoryRepo repository.ITopupHistoryRepository
}

func NewTopupHistoryUsecase(topupHistoryRepo repository.ITopupHistoryRepository) *TopupHistoryUsecase {
	return &TopupHistoryUsecase{topupHistoryRepo: topupHistoryRepo}
}

func (u *TopupHistoryUsecase) CreateTopupHistory(topupHistory *model.TopupHistory) (*model.TopupHistory, error) {
	var error []string

	if topupHistory.UserID == uuid.Nil {
		error = append(error, "user ID is required")
	}
	if topupHistory.PaymentID == "" {
		error = append(error, "payment ID is required")
	}
	if topupHistory.Amount <= 0 {
		error = append(error, "amount must be greater than 0")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.topupHistoryRepo.CreateTopupHistory(topupHistory)
}

func (u *TopupHistoryUsecase) GetTopupHistoryByID(topup_history_id int) (*model.TopupHistory, error) {
	return u.topupHistoryRepo.GetTopupHistoryByID(topup_history_id)
}

func (u *TopupHistoryUsecase) GetAllTopupHistoryByUser(userID uuid.UUID) (*[]model.TopupHistory, error) {
	return u.topupHistoryRepo.GetAllTopupHistoryByUser(userID)
}

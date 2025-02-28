package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"

	"github.com/google/uuid"
)

type TransactionUsecase struct {
	transactionRepo repository.ITransactionRepository
}

func NewTransactionUsecase(transactionRepo repository.ITransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{transactionRepo: transactionRepo}
}

func (u *TransactionUsecase) CreateTransaction(transaction *model.Transactions) (*model.Transactions, error) {
	var error []string

	if transaction.UserID == uuid.Nil {
		error = append(error, "user ID is required")
	}
	if transaction.LessorID <= 0 {
		error = append(error, "product ID is required")
	}
	if transaction.Amount <= 0 {
		error = append(error, "amount must be greater than 0")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.transactionRepo.CreateTransaction(transaction)
}

func (u *TransactionUsecase) GetTransactionByID(transactionID int) (*model.Transactions, error) {
	return u.transactionRepo.GetTransactionByID(transactionID)
}

func (u *TransactionUsecase) GetAllTransactionByUser(userID uuid.UUID) ([]model.Transactions, error) {
	return u.transactionRepo.GetAllTransactionByUser(userID)
}

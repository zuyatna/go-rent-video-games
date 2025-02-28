package repository

import (
	"rent-video-game/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	CreateTransaction(transaction *model.Transactions) (*model.Transactions, error)
	GetTransactionByID(transactionID int) (*model.Transactions, error)
	GetAllTransactionByUser(userID uuid.UUID) ([]model.Transactions, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (r *TransactionRepository) CreateTransaction(transaction *model.Transactions) (*model.Transactions, error) {
	if err := r.db.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *TransactionRepository) GetTransactionByID(transactionID int) (*model.Transactions, error) {
	var transaction model.Transactions
	if err := r.db.Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetAllTransactionByUser(userID uuid.UUID) ([]model.Transactions, error) {
	var transactions []model.Transactions
	if err := r.db.Where("user_id = ?", userID).Preload("Products").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

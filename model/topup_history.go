package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopupHistory struct {
	TopupHistoryID int            `json:"topup_history_id" gorm:"type:serial;primaryKey"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid; not null"`
	PaymentID      string         `json:"payment_id" gorm:"type:varchar(255); not null"`
	Amount         float64        `json:"amount" gorm:"type:decimal(10,2); not null"`
	CreatedAt      time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
	Users          Users          `json:"-" gorm:"foreignKey:UserID;references:UserID"`
}

type TopupHistoryRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}

type TopupHistoryData struct {
	TopupHistoryID int       `json:"topup_history_id"`
	PaymentID      string    `json:"payment_id"`
	Amount         float64   `json:"amount"`
	CreatedAt      time.Time `json:"created_at"`
}

type TopupHistoryResponse struct {
	Message string             `json:"message"`
	Data    []TopupHistoryData `json:"data"`
}

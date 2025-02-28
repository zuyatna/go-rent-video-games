package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transactions struct {
	TransactionID int            `json:"transaction_id" gorm:"type:serial;primaryKey"`
	BookingID     int            `json:"booking_id" gorm:"type:int; not null"`
	UserID        uuid.UUID      `json:"user_id" gorm:"type:uuid; not null"`
	LessorID      int            `json:"lessor_id" gorm:"type:int; not null"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2); not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
	Bookings      Bookings       `json:"-" gorm:"foreignKey:BookingID;references:BookingID"`
	Users         Users          `json:"-" gorm:"foreignKey:UserID;references:UserID"`
	Lessors       Lessors        `json:"-" gorm:"foreignKey:LessorID;references:LessorID"`
}

type TransactionRequest struct {
	BookingID int       `json:"booking_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	LessorID  int       `json:"lessor_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
}

type TransactionData struct {
	TransactionID int     `json:"transaction_id"`
	BookingID     int     `json:"booking_id"`
	ReceiveID     int     `json:"receive_id"`
	Amount        float64 `json:"amount"`
}

type TransactionResponse struct {
	Message string            `json:"message"`
	Data    []TransactionData `json:"data"`
}

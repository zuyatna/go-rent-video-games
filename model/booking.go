package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	Pending  BookingStatus = "PENDING"
	Approved BookingStatus = "APPROVED"
	Rejected BookingStatus = "REJECTED"
)

type Bookings struct {
	BookingID int            `json:"booking_id" gorm:"type:serial;primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid; not null"`
	ProductID int            `json:"product_id" gorm:"type:int; not null"`
	StartDate string         `json:"start_date" gorm:"type:date; not null"`
	EndDate   string         `json:"end_date" gorm:"type:date; not null"`
	Status    BookingStatus  `json:"status" gorm:"type:booking_status; not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
	Users     Users          `json:"-" gorm:"foreignKey:UserID;references:UserID"`
	Products  Products       `json:"-" gorm:"foreignKey:ProductID;references:ProductID"`
}

type BookingRequest struct {
	ProductID int    `json:"product_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type BookingData struct {
	BookingID   int    `json:"booking_id"`
	ProductName string `json:"product_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Status      string `json:"status"`
}

type BookingResponse struct {
	Message string        `json:"message"`
	Data    []BookingData `json:"data"`
}

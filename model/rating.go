package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ratings struct {
	RatingID  int            `json:"rating_id" gorm:"type:serial;primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid; not null"`
	ProductID int            `json:"product_id" gorm:"type:int; not null"`
	Review    string         `json:"review" gorm:"type:text"`
	Stars     float64        `json:"stars" gorm:"type:decimal(2,1); not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
}

type RatingRequest struct {
	ProductID int     `json:"product_id" validate:"required"`
	Review    string  `json:"review" validate:"required"`
	Stars     float64 `json:"stars" validate:"required"`
}

type RatingData struct {
	RatingID  int     `json:"rating_id"`
	ProductID int     `json:"product_id"`
	Review    string  `json:"review"`
	Stars     float64 `json:"stars"`
}

type RatingResponse struct {
	Message string       `json:"message"`
	Data    []RatingData `json:"data"`
}

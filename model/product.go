package model

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ProductID          int            `json:"product_id" gorm:"type:serial;primaryKey"`
	LessorID           int            `json:"lessor_id" gorm:"type:int; not null"`
	ConsoleID          int            `json:"console_id" gorm:"type:int; not null"`
	Name               string         `json:"name" gorm:"type:varchar(255); not null"`
	Description        string         `json:"description" gorm:"type:text; not null"`
	RentalCostPerMonth float64        `json:"rental_cost_per_month" gorm:"type:decimal(10,2); not null"`
	StockAvailability  int            `json:"stock_availability" gorm:"type:int; not null"`
	CreatedAt          time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
	Lessors            Lessors        `json:"-" gorm:"foreignKey:LessorID;references:LessorID"`
	Consoles           Consoles       `json:"-" gorm:"foreignKey:ConsoleID;references:ConsoleID"`
}

type ProductRequest struct {
	ConsoleID          int     `json:"console_id" validate:"required"`
	Name               string  `json:"name" validate:"required"`
	Description        string  `json:"description" validate:"required"`
	RentalCostPerMonth float64 `json:"rental_cost_per_month" validate:"required"`
	StockAvailability  int     `json:"stock_availability" validate:"required"`
}

type ProductData struct {
	ProductID          int     `json:"product_id"`
	ConsoleName        string  `json:"console_name"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	RentalCostPerMonth float64 `json:"rental_cost_per_month"`
	Stars              float64 `json:"stars"`
	StockAvailability  int     `json:"stock_availability"`
}

type ProductPublicData struct {
	ProductID          int     `json:"product_id"`
	Name               string  `json:"name"`
	RentalCostPerMonth float64 `json:"rental_cost_per_month"`
	Stars              float64 `json:"stars"`
	StockAvailability  int     `json:"stock_availability"`
	Location           string  `json:"location"`
}

type ProductResponse struct {
	Message string        `json:"message"`
	Data    []ProductData `json:"data"`
}

type ProductPublicResponse struct {
	Message string              `json:"message"`
	Data    []ProductPublicData `json:"data"`
}

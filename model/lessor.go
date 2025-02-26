package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lessors struct {
	LessorID  int            `json:"lessor_id" gorm:"type:serial;primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid; not null"`
	Name      string         `json:"name" gorm:"type:varchar(255); not null"`
	Location  string         `json:"location" gorm:"type:varchar(255); not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
	Users     Users          `json:"-" gorm:"foreignKey:UserID;references:UserID"`
}

type LessorRegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type LessorRegisterResponse struct {
	Message string `json:"message"`
	Data    struct {
		LessorID int    `json:"lessor_id"`
		Name     string `json:"name"`
		Location string `json:"location"`
	} `json:"data"`
}

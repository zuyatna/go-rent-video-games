package model

import (
	"time"

	"gorm.io/gorm"
)

type Consoles struct {
	ConsoleID int            `json:"console_id" gorm:"type:serial;primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255); not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
}

type ConsoleData struct {
	ConsoleID int    `json:"console_id"`
	Name      string `json:"name"`
}

type ConsoleResponse struct {
	Message string        `json:"message"`
	Data    []ConsoleData `json:"data"`
}

package repository

import (
	"rent-video-game/model"

	"gorm.io/gorm"
)

type IConsoleRepository interface {
	GetAllConsole() ([]model.Consoles, error)
	GetConsoleID(consoleID int) (*model.Consoles, error)
}

type ConsoleRepository struct {
	db *gorm.DB
}

func NewConsoleRepository(db *gorm.DB) *ConsoleRepository {
	return &ConsoleRepository{db}
}

func (r *ConsoleRepository) GetAllConsole() ([]model.Consoles, error) {
	var consoles []model.Consoles
	if err := r.db.Find(&consoles).Error; err != nil {
		return nil, err
	}
	return consoles, nil
}

func (r *ConsoleRepository) GetConsoleID(consoleID int) (*model.Consoles, error) {
	var console model.Consoles
	if err := r.db.Where("console_id = ?", consoleID).First(&console).Error; err != nil {
		return nil, err
	}
	return &console, nil
}

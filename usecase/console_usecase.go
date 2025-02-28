package usecase

import (
	"rent-video-game/model"
	"rent-video-game/repository"
)

type ConsoleUsecase struct {
	consoleRepo repository.IConsoleRepository
}

func NewConsoleUsecase(consoleRepo repository.IConsoleRepository) *ConsoleUsecase {
	return &ConsoleUsecase{consoleRepo: consoleRepo}
}

func (u *ConsoleUsecase) GetAllConsole() ([]model.Consoles, error) {
	return u.consoleRepo.GetAllConsole()
}

func (u *ConsoleUsecase) GetConsoleID(consoleID int) (*model.Consoles, error) {
	return u.consoleRepo.GetConsoleID(consoleID)
}

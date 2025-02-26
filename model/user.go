package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"type:varchar(255); not null"`
	Email     string         `json:"email" gorm:"type:varchar(255); not null; unique"`
	Password  string         `json:"password" gorm:"type:varchar(255); not null"`
	Address   string         `json:"address" gorm:"type:varchar(255); not null"`
	Amount    int            `json:"amount" gorm:"type:int; not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
}

func (u *Users) CompareHashAndPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID uuid.UUID `json:"user_id"`
		Email  string    `json:"email"`
		Name   string    `json:"name"`
	} `json:"data"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type TopupRequest struct {
	Amount int `json:"amount" validate:"required,gt=0"`
}

type TopupResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID      uuid.UUID `json:"user_id"`
		TopupAmount int       `json:"topup_amount"`
		NewBalance  int       `json:"new_balance"`
		PaymentID   string    `json:"payment_id"`
	} `json:"data"`
}

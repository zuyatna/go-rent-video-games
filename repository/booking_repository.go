package repository

import (
	"errors"
	"rent-video-game/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBookingRepository interface {
	CreateBooking(booking *model.Bookings) (*model.Bookings, error)
	GetBookingByID(bookingID int, userID uuid.UUID) (*model.Bookings, error)
	GetAllBookingByUser(userID uuid.UUID) ([]model.Bookings, error)
	UpdateBooking(bookingID int, status model.BookingStatus, booking *model.Bookings) (*model.Bookings, error)

	IsUserProductOwner(userID uuid.UUID, productID int) (bool, error)
	GetProductByID(productID int) (*model.Products, error)
}

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db}
}

func (r *BookingRepository) CreateBooking(booking *model.Bookings) (*model.Bookings, error) {
	if err := r.db.Create(booking).Error; err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *BookingRepository) GetBookingByID(bookingID int, userID uuid.UUID) (*model.Bookings, error) {
	var booking model.Bookings
	if err := r.db.Where("booking_id = ? AND user_id = ?", bookingID, userID).First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) GetAllBookingByUser(userID uuid.UUID) ([]model.Bookings, error) {
	var bookings []model.Bookings
	if err := r.db.Where("user_id = ?", userID).Preload("Products").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepository) UpdateBooking(bookingID int, status model.BookingStatus, booking *model.Bookings) (*model.Bookings, error) {
	var b model.Bookings
	err := r.db.Where("booking_id = ?", bookingID).Preload("Products").First(&b).Error
	if err != nil {
		return &b, err
	}

	if b.Status == model.Approved {
		return nil, errors.New("cannot update an already approved booking")
	}

	if err := r.db.Model(&model.Bookings{}).Where("booking_id = ?", bookingID).Update("status", status).Error; err != nil {
		return nil, err
	}

	if err := r.db.Where("booking_id = ?", bookingID).Preload("Products").First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) IsUserProductOwner(userID uuid.UUID, productID int) (bool, error) {
	var count int64
	err := r.db.Table("products").
		Joins("JOIN lessors ON products.lessor_id = lessors.lessor_id").
		Where("products.product_id = ? AND lessors.user_id = ? AND products.deleted_at IS NULL AND lessors.deleted_at IS NULL",
			productID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BookingRepository) GetProductByID(productID int) (*model.Products, error) {
	var product model.Products
	if err := r.db.Where("product_id = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

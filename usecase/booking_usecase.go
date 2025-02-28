package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"

	"github.com/google/uuid"
)

type BookingUsecase struct {
	bookingRepo repository.IBookingRepository
}

func NewBookingUsecase(bookingRepo repository.IBookingRepository) *BookingUsecase {
	return &BookingUsecase{bookingRepo: bookingRepo}
}

func (u *BookingUsecase) CreateBooking(booking *model.Bookings) (*model.Bookings, error) {
	var error []string

	if booking.UserID == uuid.Nil {
		error = append(error, "user ID is required")
	}
	if booking.ProductID == 0 {
		error = append(error, "product ID is required")
	}
	if booking.StartDate == "" {
		error = append(error, "start date is required")
	}
	if booking.EndDate == "" {
		error = append(error, "end date is required")
	}
	if booking.Status == "" {
		error = append(error, "status is required")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.bookingRepo.CreateBooking(booking)
}

func (u *BookingUsecase) GetBookingByID(bookingID int, userID uuid.UUID) (*model.Bookings, error) {
	return u.bookingRepo.GetBookingByID(bookingID, userID)
}

func (u *BookingUsecase) GetAllBookingByUser(userID uuid.UUID) ([]model.Bookings, error) {
	return u.bookingRepo.GetAllBookingByUser(userID)
}

func (u *BookingUsecase) UpdateBooking(bookingID int, status model.BookingStatus, booking *model.Bookings) (*model.Bookings, error) {
	currentBooking, err := u.bookingRepo.GetBookingByID(bookingID, booking.UserID)
	if err != nil {
		return nil, err
	}

	if currentBooking.Status == model.Approved {
		return nil, errors.New("cannot update an already approved booking")
	}

	return u.bookingRepo.UpdateBooking(bookingID, status, booking)
}

func (u *BookingUsecase) IsUserProductOwner(userID uuid.UUID, productID int) (bool, error) {
	return u.bookingRepo.IsUserProductOwner(userID, productID)
}

func (u *BookingUsecase) GetProductByID(productID int) (*model.Products, error) {
	return u.bookingRepo.GetProductByID(productID)
}

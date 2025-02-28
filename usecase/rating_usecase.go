package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"

	"github.com/google/uuid"
)

type RatingUsecase struct {
	ratingRepo repository.IRatingRepository
}

func NewRatingUsecase(ratingRepo repository.IRatingRepository) *RatingUsecase {
	return &RatingUsecase{ratingRepo: ratingRepo}
}

func (u *RatingUsecase) CreateRating(rating *model.Ratings) (*model.Ratings, error) {
	var error []string

	if rating.ProductID < 0 {
		error = append(error, "product ID is required")
	}
	if rating.Review == "" {
		error = append(error, "review is required")
	}
	if rating.Stars <= 0 {
		error = append(error, "stars must be greater than 0")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.ratingRepo.CreateRating(rating)
}

func (u *RatingUsecase) GetAllRatingByProduct(productID int) ([]model.Ratings, error) {
	return u.ratingRepo.GetAllRatingByProduct(productID)
}

func (u *RatingUsecase) GetAverageRatingByProduct(productID int) (float64, error) {
	return u.ratingRepo.GetAverageRatingByProduct(productID)
}

func (u *RatingUsecase) GetRatingByUserAndProduct(userID uuid.UUID, productID int) (*model.Ratings, error) {
	return u.ratingRepo.GetRatingByUserAndProduct(userID, productID)
}

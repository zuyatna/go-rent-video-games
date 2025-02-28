package repository

import (
	"rent-video-game/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRatingRepository interface {
	CreateRating(rating *model.Ratings) (*model.Ratings, error)

	GetAllRatingByProduct(productID int) ([]model.Ratings, error)
	GetAverageRatingByProduct(productID int) (float64, error)
	GetRatingByUserAndProduct(userID uuid.UUID, productID int) (*model.Ratings, error)
}

type RatingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{db}
}

func (r *RatingRepository) CreateRating(rating *model.Ratings) (*model.Ratings, error) {
	if err := r.db.Create(rating).Error; err != nil {
		return nil, err
	}
	return rating, nil
}

func (r *RatingRepository) GetAllRatingByProduct(productID int) ([]model.Ratings, error) {
	var ratings []model.Ratings
	if err := r.db.Where("product_id = ?", productID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *RatingRepository) GetAverageRatingByProduct(productID int) (float64, error) {
	var avgRating float64

	row := r.db.Model(&model.Ratings{}).
		Select("COALESCE(AVG(stars), 0) as stars").
		Where("product_id = ?", productID).
		Row()

	if err := row.Scan(&avgRating); err != nil {
		return 0, err
	}

	return avgRating, nil
}

func (r *RatingRepository) GetRatingByUserAndProduct(userID uuid.UUID, productID int) (*model.Ratings, error) {
	var rating model.Ratings
	if err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&rating).Error; err != nil {
		return nil, err
	}
	return &rating, nil
}

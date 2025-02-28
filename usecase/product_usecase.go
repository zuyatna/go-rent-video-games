package usecase

import (
	"errors"
	"rent-video-game/model"
	"rent-video-game/repository"
	"strings"
)

type ProductUsecase struct {
	productRepo repository.IProductRepository
}

func NewProductUsecase(productRepo repository.IProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) RegisterProduct(product *model.Products) (*model.Products, error) {
	var error []string

	if product.ConsoleID < 0 {
		error = append(error, "console ID is required")
	}
	if product.Name == "" {
		error = append(error, "name is required")
	}
	if product.Description == "" {
		error = append(error, "description is required")
	}
	if product.RentalCostPerMonth <= 0 {
		error = append(error, "rental cost per month must be greater than 0")
	}
	if product.StockAvailability < 0 {
		error = append(error, "stock availability must be 0 or greater")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.productRepo.RegisterProduct(product)
}

func (u *ProductUsecase) GetProductByID(productID, lessorID int) (*model.Products, error) {
	return u.productRepo.GetProductByID(productID, lessorID)
}

func (u *ProductUsecase) GetAllProductsByLessor(lessorID int) ([]model.Products, error) {
	return u.productRepo.GetAllProductsByLessor(lessorID)
}

func (u *ProductUsecase) UpdateProduct(productID int, product *model.Products) (*model.Products, error) {
	var error []string

	if product.ConsoleID < 0 {
		error = append(error, "console ID is required")
	}
	if product.Name == "" {
		error = append(error, "name is required")
	}
	if product.Description == "" {
		error = append(error, "description is required")
	}
	if product.RentalCostPerMonth <= 0 {
		error = append(error, "rental cost per month must be greater than 0")
	}
	if product.StockAvailability < 0 {
		error = append(error, "stock availability must be 0 or greater")
	}

	if len(error) > 0 {
		return nil, errors.New(strings.Join(error, ", "))
	}

	return u.productRepo.UpdateProduct(productID, product)
}

func (u *ProductUsecase) DeleteProduct(productID, lessorID int) (*model.Products, error) {
	return u.productRepo.DeleteProduct(productID, lessorID)
}

func (u *ProductUsecase) GetLessorByProductID(productID int) (*model.Lessors, error) {
	return u.productRepo.GetLessorByProductID(productID)
}

func (u *ProductUsecase) GetAllProducts() ([]model.Products, error) {
	return u.productRepo.GetAllProducts()
}

func (u *ProductUsecase) IncrementStockAvailability(productID int) error {
	return u.productRepo.IncrementStockAvailability(productID)
}

func (u *ProductUsecase) DecrementStockAvailability(productID int) error {
	return u.productRepo.DecrementStockAvailability(productID)
}

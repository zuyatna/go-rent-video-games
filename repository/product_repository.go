package repository

import (
	"rent-video-game/model"
	"time"

	"gorm.io/gorm"
)

type IProductRepository interface {
	RegisterProduct(product *model.Products) (*model.Products, error)
	GetProductByID(productID, lessorID int) (*model.Products, error)
	GetAllProductsByLessor(lessorID int) ([]model.Products, error)
	UpdateProduct(productID int, product *model.Products) (*model.Products, error)
	DeleteProduct(productID, lessorID int) (*model.Products, error)

	GetLessorByProductID(productID int) (*model.Lessors, error)
	GetAllProducts() ([]model.Products, error)

	IncrementStockAvailability(productID int) error
	DecrementStockAvailability(productID int) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) RegisterProduct(product *model.Products) (*model.Products, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) GetProductByID(productID, lessorID int) (*model.Products, error) {
	var product model.Products
	if err := r.db.Where("product_id = ? AND lessor_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		productID, lessorID, "0001-01-01 00:00:00").First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAllProductsByLessor(lessorID int) ([]model.Products, error) {
	var products []model.Products
	if err := r.db.Where("lessor_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		lessorID, "0001-01-01 00:00:00").Preload("Consoles").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) UpdateProduct(productID int, product *model.Products) (*model.Products, error) {
	var p model.Products
	err := r.db.Where("product_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		productID, "0001-01-01 00:00:00").Preload("Consoles").First(&p).Error
	if err != nil {
		return &p, err
	}

	p.ConsoleID = product.ConsoleID
	p.Name = product.Name
	p.Description = product.Description
	p.RentalCostPerMonth = product.RentalCostPerMonth
	p.StockAvailability = product.StockAvailability

	if err := r.db.Save(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) DeleteProduct(productID, lessorID int) (*model.Products, error) {
	var p model.Products
	err := r.db.Where("product_id = ? AND lessor_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		productID, lessorID, "0001-01-01 00:00:00").First(&p).Error
	if err != nil {
		return &p, err
	}

	if err := r.db.Model(&p).Update("deleted_at", time.Now()).Error; err != nil {
		return &p, err
	}
	return &p, nil
}

func (r *ProductRepository) GetLessorByProductID(productID int) (*model.Lessors, error) {
	var lessor model.Lessors
	if err := r.db.Where("product_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		productID, "0001-01-01 00:00:00").Preload("User").First(&lessor).Error; err != nil {
		return nil, err
	}
	return &lessor, nil
}

func (r *ProductRepository) GetAllProducts() ([]model.Products, error) {
	var products []model.Products
	if err := r.db.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Preload("Lessors").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) IncrementStockAvailability(productID int) error {
	return r.db.Model(&model.Products{}).Where("product_id = ?", productID).Update("stock_availability", gorm.Expr("stock_availability + 1")).Error
}

func (r *ProductRepository) DecrementStockAvailability(productID int) error {
	return r.db.Model(&model.Products{}).Where("product_id = ?", productID).Update("stock_availability", gorm.Expr("stock_availability - 1")).Error
}

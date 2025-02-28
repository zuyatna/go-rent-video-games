package handler

import (
	"net/http"
	"rent-video-game/middleware"
	"rent-video-game/model"
	"rent-video-game/usecase"
	"rent-video-game/utils"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	lessorUsecase  *usecase.LessorUsecase
	ratingUsecase  *usecase.RatingUsecase
}

func NewProductHandler(
	productUsecase *usecase.ProductUsecase,
	lessorUsecase *usecase.LessorUsecase,
	ratingUsecase *usecase.RatingUsecase,
) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		lessorUsecase:  lessorUsecase,
		ratingUsecase:  ratingUsecase,
	}
}

func (u *ProductHandler) ProductRoutes(e *echo.Echo) {
	e.POST("/lessor/product", middleware.UserAuthMiddleware()(u.RegisterProduct))
	e.GET("/lessor/product/:product_id", middleware.UserAuthMiddleware()(u.GetProductByID))
	e.GET("/lessor/products", middleware.UserAuthMiddleware()(u.GetAllProductsByLessor))
	e.PUT("/lessor/product/:product_id", middleware.UserAuthMiddleware()(u.UpdateProduct))
	e.DELETE("/lessor/product/:product_id", middleware.UserAuthMiddleware()(u.DeleteProduct))

	e.GET("/products", u.GetAllProducts)
}

func (u *ProductHandler) RegisterProduct(c echo.Context) error {
	var productReq *model.ProductRequest
	if err := c.Bind(&productReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if lessor.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden access")
	}

	product := &model.Products{
		ConsoleID:          productReq.ConsoleID,
		Name:               productReq.Name,
		Description:        productReq.Description,
		RentalCostPerMonth: productReq.RentalCostPerMonth,
		StockAvailability:  productReq.StockAvailability,
	}

	product.LessorID = lessor.LessorID // set lessor id

	product, err = u.productUsecase.RegisterProduct(product)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	productData := model.ProductData{
		ProductID:          product.ProductID,
		ConsoleName:        product.Consoles.Name,
		Name:               product.Name,
		Description:        product.Description,
		RentalCostPerMonth: product.RentalCostPerMonth,
		StockAvailability:  product.StockAvailability,
	}

	response := model.ProductResponse{
		Message: "success register product",
		Data:    []model.ProductData{productData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *ProductHandler) GetProductByID(c echo.Context) error {
	productID := c.Param("product_id")
	id := utils.StringToInt(productID)

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	product, err := u.productUsecase.GetProductByID(id, lessor.LessorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	productData := model.ProductData{
		ProductID:          product.ProductID,
		Name:               product.Name,
		Description:        product.Description,
		RentalCostPerMonth: product.RentalCostPerMonth,
		StockAvailability:  product.StockAvailability,
	}

	response := model.ProductResponse{
		Message: "success get product by id",
		Data:    []model.ProductData{productData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *ProductHandler) GetAllProductsByLessor(c echo.Context) error {
	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	products, err := u.productUsecase.GetAllProductsByLessor(lessor.LessorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var productData []model.ProductData
	for _, value := range products {

		stars, err := u.ratingUsecase.GetAverageRatingByProduct(value.ProductID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		productData = append(productData, model.ProductData{
			ProductID:          value.ProductID,
			ConsoleName:        value.Consoles.Name,
			Name:               value.Name,
			Description:        value.Description,
			RentalCostPerMonth: value.RentalCostPerMonth,
			Stars:              stars,
			StockAvailability:  value.StockAvailability,
		})
	}

	response := model.ProductResponse{
		Message: "success get all products",
		Data:    productData,
	}

	return c.JSON(http.StatusOK, response)
}

func (u *ProductHandler) UpdateProduct(c echo.Context) error {
	var product *model.Products
	if err := c.Bind(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	productID := c.Param("product_id")
	id := utils.StringToInt(productID)

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if lessor.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden access")
	}

	product.LessorID = lessor.LessorID
	product, err = u.productUsecase.UpdateProduct(id, product)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	productData := model.ProductData{
		ProductID:          product.ProductID,
		ConsoleName:        product.Consoles.Name,
		Name:               product.Name,
		Description:        product.Description,
		RentalCostPerMonth: product.RentalCostPerMonth,
		StockAvailability:  product.StockAvailability,
	}

	response := model.ProductResponse{
		Message: "success update product",
		Data:    []model.ProductData{productData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *ProductHandler) DeleteProduct(c echo.Context) error {
	productID := c.Param("product_id")
	id := utils.StringToInt(productID)

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	product, err := u.productUsecase.DeleteProduct(id, lessor.LessorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	productData := model.ProductData{
		ProductID:          product.ProductID,
		Name:               product.Name,
		Description:        product.Description,
		RentalCostPerMonth: product.RentalCostPerMonth,
		StockAvailability:  product.StockAvailability,
	}

	response := model.ProductResponse{
		Message: "success delete product",
		Data:    []model.ProductData{productData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := u.productUsecase.GetAllProducts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var productData []model.ProductPublicData
	for _, value := range products {

		stars, err := u.ratingUsecase.GetAverageRatingByProduct(value.ProductID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		productData = append(productData, model.ProductPublicData{
			ProductID:          value.ProductID,
			Name:               value.Name,
			RentalCostPerMonth: value.RentalCostPerMonth,
			Stars:              stars,
			StockAvailability:  value.StockAvailability,
			Location:           value.Lessors.Location,
		})
	}

	response := model.ProductPublicResponse{
		Message: "success get all products",
		Data:    productData,
	}

	return c.JSON(http.StatusOK, response)
}

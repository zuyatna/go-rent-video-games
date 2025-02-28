package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rent-video-game/config"
	"rent-video-game/handler"
	"rent-video-game/model"
	"rent-video-game/repository"
	"rent-video-game/usecase"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// create shutdown signal listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db, err := gorm.Open(postgres.Open(config.InitDB()), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database!")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to connect to database!")
	}

	// sceheme migration
	db.AutoMigrate(
		&model.Users{},
		&model.Lessors{},
		&model.Consoles{},
		&model.TopupHistory{},
		&model.Products{},
		&model.Bookings{},
		&model.Transactions{},
		&model.Ratings{},
	)
	fmt.Println("database migrated")

	// init echo
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// topup history handler
	topupHistoryRepo := repository.NewTopupHistoryRepository(db)
	topupHistoryUsecase := usecase.NewTopupHistoryUsecase(topupHistoryRepo)
	topupHistoryHandler := handler.NewTopupHistoryHandler(topupHistoryUsecase)
	topupHistoryHandler.TopupHistoryRoutes(e)

	// user handler
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase, topupHistoryUsecase)
	userHandler.UserRoutes(e)

	// lessor handler
	lessorRepo := repository.NewLessorRepository(db)
	lessorUsecase := usecase.NewLessorUsecase(lessorRepo)
	lessorHandler := handler.NewLessorHandler(lessorUsecase)
	lessorHandler.LessorRoutes(e)

	// console handler
	consoleRepo := repository.NewConsoleRepository(db)
	consoleUsecase := usecase.NewConsoleUsecase(consoleRepo)
	consoleHandler := handler.NewConsoleHandler(consoleUsecase)
	consoleHandler.ConsoleRoutes(e)

	// rating handler
	ratingRepo := repository.NewRatingRepository(db)
	ratingUsecase := usecase.NewRatingUsecase(ratingRepo)
	ratingHandler := handler.NewRatingHandler(ratingUsecase)
	ratingHandler.RatingRoutes(e)

	// product handler
	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase, lessorUsecase, ratingUsecase)
	productHandler.ProductRoutes(e)

	// booking handler
	bookingRepo := repository.NewBookingRepository(db)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingUsecase, userUsecase, productUsecase, lessorUsecase)
	bookingHandler.BookingRoutes(e)

	// transaction handler
	transactionRepo := repository.NewTransactionRepository(db)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase, bookingUsecase, userUsecase, lessorUsecase)
	transactionHandler.TransactionRoutes(e)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	// waiting for shutdown signal
	<-quit
	fmt.Println("shutting down server...")

	// give server 10 seconds to finish processing requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("server forced to shutdown:", err)
	}

	if err := sqlDB.Close(); err != nil {
		e.Logger.Fatal("failed to close database connection", err)
	}

	fmt.Println("server exited properly")
}

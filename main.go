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
	)
	fmt.Println("database migrated")

	// init echo
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/success", func(c echo.Context) error {
		return c.String(http.StatusOK, "Success!")
	})

	// user handler
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	userHandler.UserRoutes(e)

	// lessor handler
	lessorRepo := repository.NewLessorRepository(db)
	lessorUsecase := usecase.NewLessorUsecase(lessorRepo)
	lessorHandler := handler.NewLessorHandler(lessorUsecase)
	lessorHandler.LessorRoutes(e)

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

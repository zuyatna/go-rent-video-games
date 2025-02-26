package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func InitDB() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}

	return os.Getenv("DATABASE_URL")
}

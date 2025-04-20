package main

import (
	"log"

	"github.com/drezza544/go-apibank/initializers"
	"github.com/drezza544/go-apibank/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectionDatabase()
}

func main() {
	err := initializers.DB.AutoMigrate(&models.User{}, &models.Bank{}, &models.Account{}, &models.Transaction{})

	if err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	log.Println("✅ Auto migration complete")
}

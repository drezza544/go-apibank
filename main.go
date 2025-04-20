package main

import (
	"fmt"

	"github.com/drezza544/go-apibank/controllers"
	"github.com/drezza544/go-apibank/initializers"
	"github.com/drezza544/go-apibank/routes"
	"github.com/drezza544/go-apibank/services"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectionDatabase()
}

func main() {
	fmt.Println("Hello World")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	userService := services.NewUserService(initializers.DB)
	userController := controllers.NewUserController(userService)

	bankService := services.NewBankService(initializers.DB)
	bankController := controllers.NewBankController(bankService)

	transactionService := services.NewTransactionService(initializers.DB)
	transactionController := controllers.NewTransactionController(transactionService)

	routes.RegisterUserRoutes(r, userController)
	routes.RegisterBankRoutes(r, bankController)
	routes.RegisterTransactionRoutes(r, transactionController)

	// r.GET("/users", controllers.AllUsers)
	// // r.GET("/users/:id", controllers.GetUsers)
	// r.POST("/users", controllers.NewUserController(initializers.DB))
	// // r.PUT("/users/:id", controllers.UpdateUsers)

	// r.GET("/bank", controllers.AllBanks)
	// r.GET("/bank/:id", controllers.AllBanks)
	// r.POST("/bank", controllers.CreateBank)
	// r.PUT("/bank/:id", controllers.UpdateBank)

	r.Run()
}

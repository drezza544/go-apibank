// routes/routes.go
package routes

import (
	"github.com/drezza544/go-apibank/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userController *controllers.UserController) {
	r.GET("/users", userController.AllUsers)
	r.POST("/users", userController.CreateUsers)
	r.GET("/users/:id", userController.GetUsers)
	r.PATCH("/users/:id", userController.UpdateUsers)
}

func RegisterBankRoutes(r *gin.Engine, bankController *controllers.BankController) {
	r.GET("/bank", bankController.AllBanks)
	r.GET("/bank/:id", bankController.GetBankId)
	r.GET("/bank/code/:code", bankController.GetBankCode)
	r.POST("/bank", bankController.CreateBanks)
	r.PATCH("/bank/:id", bankController.UpdateBanks)
}

func RegisterTransactionRoutes(r *gin.Engine, transactionController *controllers.TransactionController) {
	r.GET("/transaction", transactionController.AllTransactions)
	r.POST("/transaction", transactionController.CreateTransaction)
	// r.GET("/transaction/:id", transactionController.GetTransactions)
	// r.PATCH("/transaction/:id", transactionController.UpdateTransactions)
}

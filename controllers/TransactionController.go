package controllers

import (
	"net/http"

	"github.com/drezza544/go-apibank/dto"
	"github.com/drezza544/go-apibank/services"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService services.TransactionService
}

func NewTransactionController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: *transactionService,
	}
}

func (uc *TransactionController) CreateTransaction(c *gin.Context) {
	var transactionRequest dto.TransactionRequest
	if err := c.ShouldBindJSON(&transactionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil service untuk membuat transaction
	transactions, err := uc.TransactionService.CreateTransaction(transactionRequest)
	if err != nil {
		// Jika error terjadi, kembalikan response error dengan status 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika berhasil, kembalikan transaction dengan status 201 (Created)
	c.JSON(http.StatusCreated, gin.H{"transaction": transactions})
}

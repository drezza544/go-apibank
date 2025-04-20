package controllers

import (
	"net/http"

	"github.com/drezza544/go-apibank/dto"
	"github.com/drezza544/go-apibank/services"
	"github.com/gin-gonic/gin"
)

type BankController struct {
	BankService services.BankService
}

func NewBankController(bankService *services.BankService) *BankController {
	return &BankController{
		BankService: *bankService,
	}
}

func (uc *BankController) AllBanks(c *gin.Context) {
	banks, err := uc.BankService.GetAllBanks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"data":    banks,
	})
}

func (uc *BankController) GetBankId(c *gin.Context) {
	id := c.Param("id")
	banks, err := uc.BankService.GetBankById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"data":    banks,
	})
}

func (uc *BankController) GetBankCode(c *gin.Context) {
	code := c.Param("code")
	banks, err := uc.BankService.GetBankByCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"data":    banks,
	})
}

func (uc *BankController) CreateBanks(c *gin.Context) {
	var bankRequest dto.BankRequest
	if err := c.ShouldBindJSON(&bankRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil service untuk membuat bank
	banks, err := uc.BankService.CreateBank(bankRequest)
	if err != nil {
		// Jika error terjadi, kembalikan response error dengan status 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika berhasil, kembalikan user dengan status 201 (Created)
	c.JSON(http.StatusCreated, gin.H{"bank": banks})
}

func (uc *BankController) UpdateBanks(c *gin.Context) {
	id := c.Param("id")
	var bankRequest dto.BankRequest
	if err := c.ShouldBindJSON(&bankRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil service untuk memperbarui banks
	banks, err := uc.BankService.UpdateBank(id, bankRequest)
	if err != nil {
		// Jika error terjadi, kembalikan response error dengan status 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika berhasil, kembalikan banks dengan status 200 (OK)
	c.JSON(http.StatusOK, gin.H{"user": banks})
}

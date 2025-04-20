// package controllers

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/drezza544/go-apibank/initializers"
// 	"github.com/drezza544/go-apibank/models"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type UserReq struct {
// 	Name     string     `json:"name" binding:"required"`
// 	Email    string     `json:"email" binding:"required"`
// 	Password string     `json:"password" binding:"required"`
// 	Type     string     `json:"type" binding:"required"`
// 	Phone    string     `json:"phone" binding:"required"`
// 	Accounts AccountReq `json:"accounts" binding:"required"`
// }

// type AccountReq struct {
// 	UserId        uuid.UUID `json:"user_id"`
// 	BankId        uuid.UUID `json:"bank_id"`
// 	AccountNumber string    `gorm:"unique" json:"account_number" binding:"required"`
// 	Balance       int64     `json:"balance" binding:"required"`
// 	Status        string    `json:"status" binding:"required"`
// }

// func AllUsers(c *gin.Context) {
// 	var users = []models.User{}
// 	if err := initializers.DB.Preload("Accounts").Preload("Accounts.Bank").Find(&users).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": users})
// }

// func GetUsers(c *gin.Context) {
// 	id := c.Param("id")
// 	var users = models.User{}
// 	if err := initializers.DB.Find(&users, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found..."})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": users})
// }

// func CreateUsers(c *gin.Context) {
// 	var users UserReq
// 	if err := c.ShouldBindJSON(&users); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	log.Printf("Parsed Users: %+v", users)

// 	if users.Accounts.BankId == uuid.Nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "BankId is required"})
// 		return
// 	}

// 	fmt.Println(users.Accounts.BankId)

// 	// 🔍 3. Validasi apakah BankId valid (ada di database)
// 	var bank models.Bank
// 	log.Println("🚨 Before querying banks")
// 	if err := initializers.DB.Debug().Where("id = ?", users.Accounts.BankId).First(&bank).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "BankId not found"})
// 		return
// 	}
// 	log.Println("✅ After querying banks")

// 	// ✅ Mulai transaksi
// 	tx := initializers.DB.Begin()

// 	user := models.User{
// 		Name:     users.Name,
// 		Email:    users.Email,
// 		Password: users.Password,
// 		Type:     users.Type,
// 		Phone:    users.Phone,
// 	}
// 	if err := tx.Create(&user).Error; err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// ✅ Buat account dengan UserId yang baru dibuat
// 	account := models.Account{
// 		UserId:        user.ID,
// 		BankId:        users.Accounts.BankId,
// 		AccountNumber: users.Accounts.AccountNumber,
// 		Balance:       users.Accounts.Balance,
// 		Status:        users.Accounts.Status,
// 	}
// 	if err := tx.Create(&account).Error; err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// ✅ Commit transaksi jika semua sukses
// 	if err := tx.Commit().Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
// 		return
// 	}

// 	log.Printf("🧪 Accounts: %+v\n", users.Accounts)
// 	log.Printf("🧪 Banks (after query): %+v\n", bank)

// 	c.JSON(http.StatusCreated, gin.H{"data": users})
// }

// func UpdateUsers(c *gin.Context) {
// 	id := c.Param("id")
// 	var users UserReq
// 	// var body struct {
// 	// 	Name     string `json:"name" binding:"required"`
// 	// 	Email    string `json:"email" binding:"required"`
// 	// 	Password string `json:"password" binding:"required"`
// 	// 	Type     string `json:"type" binding:"required"`
// 	// 	Phone    string `json:"phone" binding:"required"`
// 	// }

// 	if err := c.ShouldBindJSON(&users); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user models.User
// 	if err := initializers.DB.First(&user, id).Error; err != nil { // Check if the user exists
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}

// 	user.Name = users.Name
// 	user.Email = users.Email
// 	user.Password = users.Password
// 	user.Type = users.Type
// 	user.Phone = users.Phone

// 	user.Accounts.BankId = users.Accounts.BankId
// 	user.Accounts.AccountNumber = users.Accounts.AccountNumber
// 	user.Accounts.Balance = users.Accounts.Balance
// 	user.Accounts.Status = users.Accounts.Status

// 	if err := initializers.DB.Save(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := initializers.DB.Save(&user.Accounts).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": user})
// }

package controllers

import (
	"net/http"

	"github.com/drezza544/go-apibank/dto"
	"github.com/drezza544/go-apibank/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: *userService,
	}
}

// GET /users
func (uc *UserController) AllUsers(c *gin.Context) {
	users, err := uc.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"data":    users,
	})
}

// GET /users/:id
func (uc *UserController) GetUsers(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.UserService.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"data":    user,
	})
}

// POST /users
func (uc *UserController) CreateUsers(c *gin.Context) {
	var userRequest dto.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil service untuk membuat user
	user, err := uc.UserService.CreateUser(userRequest)
	if err != nil {
		// Jika error terjadi, kembalikan response error dengan status 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika berhasil, kembalikan user dengan status 201 (Created)
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// PATCH /users/:id
func (uc *UserController) UpdateUsers(c *gin.Context) {
	id := c.Param("id")
	var userRequest dto.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil service untuk memperbarui user
	user, err := uc.UserService.UpdateUser(id, userRequest)
	if err != nil {
		// Jika error terjadi, kembalikan response error dengan status 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika berhasil, kembalikan user dengan status 200 (OK)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

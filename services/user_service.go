package services

import (
	"fmt"
	"log"

	"github.com/drezza544/go-apibank/dto"
	"github.com/drezza544/go-apibank/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	var users []models.User

	if err := s.DB.Preload("Accounts").Preload("Accounts.Bank").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data users: %w", err)
	}

	// Mapping ke DTO
	var userResponses []dto.UserResponse
	for _, user := range users {
		var account *dto.AccountResponse
		if user.Accounts.ID != uuid.Nil {
			account = &dto.AccountResponse{
				BankId:        user.Accounts.BankId,
				AccountNumber: user.Accounts.AccountNumber,
				Balance:       user.Accounts.Balance,
				Status:        user.Accounts.Status,
			}
		}
		userReponse := dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Type:     user.Type,
			Accounts: account,
		}
		userResponses = append(userResponses, userReponse)
	}

	return userResponses, nil
}
func (s *UserService) GetUserById(id string) (*dto.UserResponse, error) {
	var user models.User

	if err := s.DB.Preload("Accounts").Preload("Accounts.Bank").First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Mapping ke DTO
	var account *dto.AccountResponse
	if user.Accounts.ID != uuid.Nil {
		account = &dto.AccountResponse{
			BankId:        user.Accounts.BankId,
			AccountNumber: user.Accounts.AccountNumber,
			Balance:       user.Accounts.Balance,
			Status:        user.Accounts.Status,
		}
	}
	userReponse := &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Type:     user.Type,
		Accounts: account,
	}

	return userReponse, nil
}

func (s *UserService) CreateUser(userRequest dto.UserRequest) (*models.User, error) {
	// Required Field
	if userRequest.Name == nil || *userRequest.Name == "" ||
		userRequest.Email == nil || *userRequest.Email == "" ||
		userRequest.Password == nil || *userRequest.Password == "" ||
		userRequest.Type == nil || *userRequest.Type == "" ||
		userRequest.Phone == nil || *userRequest.Phone == "" ||
		userRequest.Accounts == nil ||
		userRequest.Accounts.BankId == nil || *userRequest.Accounts.BankId == uuid.Nil ||
		userRequest.Accounts.AccountNumber == nil || *userRequest.Accounts.AccountNumber == "" ||
		userRequest.Accounts.Balance == nil || *userRequest.Accounts.Balance == 0 {
		return nil, fmt.Errorf("semua field wajib diisi")
	}

	// Validasi apakah nomor telepon sudah terdaftar
	var existingUser models.User
	err := s.DB.Where("phone = ? OR email = ?", userRequest.Phone, userRequest.Email).First(&existingUser).Error
	if err == nil {
		if existingUser.Phone == *userRequest.Phone {
			return nil, fmt.Errorf("nomor telepon %s sudah terdaftar", *userRequest.Phone)
		}

		return nil, fmt.Errorf("email %s sudah terdaftar", *userRequest.Email)
	}

	// Validasi apakah bank sudah terdaftar
	var existingBank models.Bank
	err = s.DB.Where("id = ?", userRequest.Accounts.BankId).First(&existingBank).Error
	if err != nil {
		return nil, fmt.Errorf("bank dengan ID %s tidak ditemukan: %w", *userRequest.Accounts.BankId, err)
	}

	// Validasi apakah nomor rekening sudah terdaftar
	var existingAccount models.Account
	err = s.DB.First(&existingAccount, "account_number = ?", userRequest.Accounts.AccountNumber).Error
	if err == nil {
		return nil, fmt.Errorf("nomor rekening %s sudah terdaftar", *userRequest.Accounts.AccountNumber)
	}

	// 💾 Simpan user dan account
	user := models.User{
		Name:     *userRequest.Name,
		Email:    *userRequest.Email,
		Password: *userRequest.Password,
		Type:     *userRequest.Type,
		Phone:    *userRequest.Phone,
		Accounts: models.Account{
			AccountNumber: *userRequest.Accounts.AccountNumber,
			BankId:        *userRequest.Accounts.BankId,
			Balance:       *userRequest.Accounts.Balance,
			Status:        "active",
		},
	}

	// Mulai transaksi
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("gagal mulai transaksi: %w", tx.Error)
	}

	// Simpan user dan account ke database dalam transaksi
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("❌ Error saat menyimpan user: %v", err)
		return nil, fmt.Errorf("gagal menyimpan user: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("❌ Error saat commit transaksi: %v", err)
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	// Jika berhasil, kembalikan user
	log.Printf("✅ User berhasil dibuat: %v", user)
	return &user, nil
}

func (s *UserService) UpdateUser(id string, userRequest dto.UserRequest) (*models.User, error) {
	var existingUser models.User

	// Validasi apakah user dengan ID tersebut ada di database
	if err := s.DB.Preload("Accounts").First(&existingUser, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user dengan ID %s tidak ditemukan: %w", id, err)
	}

	// Validasi apakah nomor telepon sudah terdaftar
	if userRequest.Phone != nil && *userRequest.Phone != "" && *userRequest.Phone != existingUser.Phone {
		if err := s.DB.Where("phone = ?", userRequest.Phone).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("nomor telepon %s sudah terdaftar", *userRequest.Phone)
		}

		existingUser.Phone = *userRequest.Phone
	}

	// Validasi apakah email sudah terdaftar
	if userRequest.Email != nil && *userRequest.Email != "" && *userRequest.Email != existingUser.Email {
		if err := s.DB.Where("email = ?", userRequest.Email).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("email %s sudah terdaftar", *userRequest.Email)
		}

		existingUser.Email = *userRequest.Email
	}

	// Update field lain jika dikirim
	if userRequest.Name != nil && *userRequest.Name != "" {
		existingUser.Name = *userRequest.Name
	}

	if userRequest.Password != nil && *userRequest.Password != "" {
		existingUser.Password = *userRequest.Password
	}

	if userRequest.Type != nil && *userRequest.Type != "" {
		existingUser.Type = *userRequest.Type
	}

	// Validasi account number
	if (userRequest.Accounts != nil && userRequest.Accounts.AccountNumber != nil) && *userRequest.Accounts.AccountNumber != "" && *userRequest.Accounts.AccountNumber != existingUser.Accounts.AccountNumber {
		var existingAccount models.Account
		if err := s.DB.Where("account_number = ?", userRequest.Accounts.AccountNumber).First(&existingAccount).Error; err == nil {
			return nil, fmt.Errorf("nomor rekening %s sudah terdaftar", *userRequest.Accounts.AccountNumber)
		}

		fmt.Println("Incoming:", *userRequest.Accounts.AccountNumber)
		fmt.Println("Existing:", existingUser.Accounts.AccountNumber)
		existingUser.Accounts.AccountNumber = *userRequest.Accounts.AccountNumber
	}

	if userRequest.Accounts != nil && userRequest.Accounts.Balance != nil && *userRequest.Accounts.Balance != 0 {
		existingUser.Accounts.Balance = *userRequest.Accounts.Balance
	}

	if userRequest.Accounts != nil && userRequest.Accounts.Status != nil && *userRequest.Accounts.Status != "" {
		existingUser.Accounts.Status = *userRequest.Accounts.Status
	}

	// Mulai transaksi
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("gagal mulai transaksi: %w", tx.Error)
	}

	// Simpan user dan account ke database dalam transaksi
	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existingUser).Error; err != nil {
		tx.Rollback()
		log.Printf("❌ Error saat menyimpan user: %v", err)
		return nil, fmt.Errorf("gagal menyimpan user: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("❌ Error saat commit transaksi: %v", err)
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	// Jika berhasil, kembalikan user
	log.Printf("✅ User berhasil diupdate: %v", userRequest)
	return &existingUser, nil
}

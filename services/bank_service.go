package services

import (
	"fmt"

	"github.com/drezza544/go-apibank/dto"
	"github.com/drezza544/go-apibank/models"
	"gorm.io/gorm"
)

type BankService struct {
	DB *gorm.DB
}

func NewBankService(db *gorm.DB) *BankService {
	return &BankService{
		DB: db,
	}
}

func (s *BankService) GetAllBanks() ([]models.Bank, error) {
	var banks []models.Bank
	if err := s.DB.Find(&banks).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data banks: %w", err)
	}

	return banks, nil
}

func (s *BankService) GetBankById(id string) (*models.Bank, error) {
	var bank models.Bank
	if err := s.DB.Preload("Accounts").First(&bank, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data bank: %w", err)
	}

	return &bank, nil
}

func (s *BankService) GetBankByCode(code string) (*models.Bank, error) {
	var bank models.Bank
	if err := s.DB.Where("code = ?", code).First(&bank).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data bank: %w", err)
	}

	return &bank, nil
}

func (s *BankService) CreateBank(bankRequest dto.BankRequest) (*models.Bank, error) {
	// Field required
	if bankRequest.Code == nil || *bankRequest.Code == "" ||
		bankRequest.Name == nil || *bankRequest.Name == "" ||
		bankRequest.Address == nil || *bankRequest.Address == "" ||
		bankRequest.CostTransaction == nil || *bankRequest.CostTransaction < 0 {
		return nil, fmt.Errorf("semua field wajib diisi")
	}

	// Validasi apakah kode bank sudah terdaftar
	var existingBank models.Bank
	err := s.DB.Where("code = ?", bankRequest.Code).First(&existingBank).Error
	fmt.Println(err)
	if err == nil {
		return nil, fmt.Errorf("kode bank %s sudah terdaftar", *bankRequest.Code)
	}

	banks := models.Bank{
		Code:            *bankRequest.Code,
		Name:            *bankRequest.Name,
		Address:         *bankRequest.Address,
		CostTransaction: *bankRequest.CostTransaction,
	}

	// Mulai transaksi
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("gagal mulai transaksi: %w", tx.Error)
	}

	// 💾 Simpan bank
	if err := tx.Create(&banks).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal menyimpan bank: %w", err)
	}

	// Commit Transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return &banks, nil
}

func (s *BankService) UpdateBank(id string, bankRequest dto.BankRequest) (*models.Bank, error) {
	// Validasi apakah kode bank sudah terdaftar
	var existingBank models.Bank

	// Validasi apakah bank dengan ID tersebut ada di database
	if err := s.DB.First(&existingBank, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("bank dengan ID %s tidak ditemukan: %w", id, err)
	}

	// 💾 Simpan bank
	// Update bank
	if bankRequest.Code != nil && *bankRequest.Code != "" && *bankRequest.Code != existingBank.Code {
		// Validasi apakah kode bank sudah terdaftar
		if err := s.DB.Where("code = ?", bankRequest.Code).First(&existingBank).Error; err == nil {
			return nil, fmt.Errorf("kode bank %s sudah terdaftar", *bankRequest.Code)
		}

		existingBank.Code = *bankRequest.Code
	}
	if bankRequest.Name != nil && *bankRequest.Name != "" {
		existingBank.Name = *bankRequest.Name
	}
	if bankRequest.Address != nil && *bankRequest.Address != "" {
		existingBank.Address = *bankRequest.Address
	}
	if bankRequest.CostTransaction != nil && *bankRequest.CostTransaction != 0 {
		existingBank.CostTransaction = *bankRequest.CostTransaction
	}

	// Mulai transaksi
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("gagal mulai transaksi: %w", tx.Error)
	}

	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existingBank).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal menyimpan bank: %w", err)
	}

	// Commit Transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return &existingBank, nil
}

package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/drezza544/go-apibank/dto"
	"github.com/drezza544/go-apibank/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	DB *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		DB: db,
	}
}

func (s *TransactionService) GetAllTransaction() ([]dto.TransactionResponse, error) {
	var transactions []models.Transaction

	// Ambil transaksi dan preload relasinya
	if err := s.DB.Preload("User").Preload("User.Accounts").Preload("User.Accounts.Bank").Preload("Bank").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data transaksi: %w", err)
	}

	// Mapping ke DTO
	var transactionResponses []dto.TransactionResponse
	for _, transaction := range transactions {
		// Ambil user penerima berdasarkan nomor rekening tujuan
		var toUser models.User
		var toUserResponse *dto.ToUserResponse = nil
		err := s.DB.Preload("Accounts").Preload("Accounts.Bank").Where("id = ?", transaction.Destination_UserId).First(&toUser).Error

		if err == nil {
			toUserResponse = &dto.ToUserResponse{
				ID:   toUser.ID,
				Name: toUser.Name,
				ToAccountResponse: dto.ToAccountResponse{
					BankId:        toUser.Accounts.BankId,
					AccountNumber: toUser.Accounts.AccountNumber,
					Balance:       toUser.Accounts.Balance,
					Status:        toUser.Accounts.Status,
					ToBankResponse: dto.ToBankResponse{
						ID:   toUser.Accounts.Bank.ID,
						Name: toUser.Accounts.Bank.Name,
					},
				},
			}
		}

		fromAccount := transaction.User.Accounts
		// toAccount := toUser.Accounts

		transactionResponses = append(transactionResponses, dto.TransactionResponse{
			TransactionId:       transaction.TransactionId,
			UserId:              transaction.UserId,
			BankId:              transaction.BankId,
			Destination_Account: transaction.Destination_Account,
			Amount:              transaction.Amount,
			BalanceBefore:       transaction.BalanceBefore,
			BalanceAfter:        transaction.BalanceAfter,
			Type:                transaction.Type,
			Status:              transaction.Status,

			FromUserResponse: &dto.FromUserResponse{
				ID:   transaction.User.ID,
				Name: transaction.User.Name,
				FromAccountResponse: dto.FromAccountResponse{
					BankId:        fromAccount.BankId,
					AccountNumber: fromAccount.AccountNumber,
					Balance:       fromAccount.Balance,
					Status:        fromAccount.Status,
					FromBankResponse: dto.FromBankResponse{
						ID:   transaction.BankId,
						Name: transaction.Bank.Name,
					},
				},
			},

			ToUserResponse: toUserResponse,

			// ToUserResponse: &dto.ToUserResponse{
			// 	ID:   toUser.ID,
			// 	Name: toUser.Name,
			// 	ToAccountResponse: dto.ToAccountResponse{
			// 		BankId:        toAccount.BankId,
			// 		AccountNumber: toAccount.AccountNumber,
			// 		Balance:       toAccount.Balance,
			// 		Status:        toAccount.Status,
			// 		ToBankResponse: dto.ToBankResponse{
			// 			ID:   toAccount.Bank.ID,
			// 			Name: toAccount.Bank.Name,
			// 		},
			// 	},
			// },
		})
	}

	return transactionResponses, nil
}

func (s *TransactionService) CreateTransaction(transactionRequest dto.TransactionRequest) (*models.Transaction, error) {
	// Required Field
	if transactionRequest.UserId == nil || *transactionRequest.UserId == uuid.Nil ||
		transactionRequest.BankId == nil || *transactionRequest.BankId == uuid.Nil ||
		transactionRequest.Destination_Account == nil || *transactionRequest.Destination_Account == "" ||
		transactionRequest.Amount == nil || *transactionRequest.Amount == 0 ||
		transactionRequest.Type == nil || *transactionRequest.Type == "" ||
		transactionRequest.Status == nil || *transactionRequest.Status == "" {
		return nil, fmt.Errorf("semua field wajib diisi")
	}

	randomNumber := rand.Intn(9999999)
	generateNumber := fmt.Sprintf("TRX/%s/%06d", time.Now().Format("200601"), randomNumber)

	// Validasi apakah transaksi dengan ID tersebut sudah ada di database
	if err := s.DB.Where("transaction_id = ?", generateNumber).First(&models.Transaction{}).Error; err == nil {
		return nil, fmt.Errorf("transaksi dengan ID %s sudah ada di database", generateNumber)
	}

	// Validasi apakah user terdaftar di database
	if err := s.DB.Where("id = ?", *transactionRequest.UserId).First(&models.User{}).Error; err != nil {
		return nil, fmt.Errorf("user dengan ID %s tidak ditemukan: %w", *transactionRequest.UserId, err)
	}

	// Validasi apakah bank terdaftar di database
	if err := s.DB.Where("id = ?", *transactionRequest.BankId).First(&models.Bank{}).Error; err != nil {
		return nil, fmt.Errorf("bank dengan ID %s tidak ditemukan: %w", *transactionRequest.BankId, err)
	}

	// Validasi apakah destination account terdaftar di database
	if err := s.DB.Where("account_number = ?", *transactionRequest.Destination_Account).First(&models.Account{}).Error; err != nil {
		return nil, fmt.Errorf("destination account %s tidak ditemukan: %w", *transactionRequest.Destination_Account, err)
	}

	// Validasi minimal 10.000 transaction
	if *transactionRequest.Amount < 10000 {
		return nil, fmt.Errorf("minimal transaksi adalah 10.000")
	}

	// Mengambil data users
	var user models.User
	if err := s.DB.Preload("Accounts").First(&user, "id = ?", *transactionRequest.UserId).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Validasi apakah nomer rekening tidak boleh mengirim ke dirinya sendiri
	if user.Accounts.AccountNumber == *transactionRequest.Destination_Account {
		return nil, fmt.Errorf("tidak bisa transfer, ke rekening sendiri")
	}

	if user.Accounts.Balance < *transactionRequest.Amount {
		return nil, fmt.Errorf("saldo tidak mencukupi")
	}

	// Mulai Transaksi
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("gagal mulai transaksi: %w", tx.Error)
	}

	balanceBefore := user.Accounts.Balance
	balanceAfter := balanceBefore - *transactionRequest.Amount

	// Update saldo di akun user
	if err := tx.Model(&user.Accounts).Update("balance", balanceAfter).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal mengupdate saldo user: %w", err)
	}

	// Mengambil data user penerima
	var receiverUser models.Account
	if err := s.DB.Where("account_number = ?", *transactionRequest.Destination_Account).First(&receiverUser).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data user penerima: %w", err)
	}

	// Validasi apakah nomer rekening harus dibank yang terdaftar
	if receiverUser.BankId != *transactionRequest.BankId {
		return nil, fmt.Errorf("nomer rekening tidak terdaftar dibank")
	}

	// Update saldo di akun user penerima
	addingAmountReceiver := receiverUser.Balance + *transactionRequest.Amount
	if err := tx.Model(&receiverUser).Update("balance", addingAmountReceiver).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal mengupdate saldo user penerima: %w", err)
	}

	transaction := models.Transaction{
		UserId:              *transactionRequest.UserId,
		BankId:              *transactionRequest.BankId,
		Destination_UserId:  receiverUser.UserId,
		From_BankId:         user.Accounts.BankId,
		TransactionId:       generateNumber,
		Destination_Account: *transactionRequest.Destination_Account,
		Amount:              *transactionRequest.Amount,
		BalanceBefore:       balanceBefore,
		BalanceAfter:        balanceAfter,
		Type:                *transactionRequest.Type,
		Status:              *transactionRequest.Status,
	}
	// Simpan transaksi ke database
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal menyimpan transaksi: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return &transaction, nil
}

package dto

import "github.com/google/uuid"

type TransactionRequest struct {
	UserId              *uuid.UUID `json:"user_id"`
	BankId              *uuid.UUID `json:"bank_id"`
	Destination_UserId  *uuid.UUID `json:"destination_user_id"`
	From_BankId         *uuid.UUID `json:"from_bank_id"`
	TransactionId       *string    `json:"transaction_id"`
	Destination_Account *string    `json:"destination_account"`
	Amount              *int64     `json:"amount"`
	BalanceBefore       *int64     `json:"balance_before"`
	BalanceAfter        *int64     `json:"balance_after"`
	Type                *string    `json:"type"`
	Status              *string    `json:"status"`
}

type TransactionResponse struct {
	UserId              uuid.UUID `json:"user_id"`
	BankId              uuid.UUID `json:"bank_id"`
	TransactionId       string    `json:"transaction_id"`
	Destination_Account string    `json:"destination_account"`
	Amount              int64     `json:"amount"`
	BalanceBefore       int64     `json:"balance_before"`
	BalanceAfter        int64     `json:"balance_after"`
	Type                string    `json:"type"`
	Status              string    `json:"status"`

	FromUserResponse *FromUserResponse
	ToUserResponse   *ToUserResponse
}

type FromUserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	FromAccountResponse FromAccountResponse
}

type FromAccountResponse struct {
	BankId        uuid.UUID `json:"bank_id"`
	AccountNumber string    `json:"account_number"`
	Balance       int64     `json:"balance"`
	Status        string    `json:"status"`

	FromBankResponse FromBankResponse
}

type FromBankResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ToUserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	ToAccountResponse ToAccountResponse
}

type ToAccountResponse struct {
	BankId        uuid.UUID `json:"bank_id"`
	AccountNumber string    `json:"account_number"`
	Balance       int64     `json:"balance"`
	Status        string    `json:"status"`

	ToBankResponse ToBankResponse
}

type ToBankResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

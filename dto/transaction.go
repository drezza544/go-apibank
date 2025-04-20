package dto

import "github.com/google/uuid"

type TransactionRequest struct {
	UserId              *uuid.UUID `json:"user_id"`
	BankId              *uuid.UUID `json:"bank_id"`
	TransactionId       *string    `json:"transaction_id"`
	Destination_Account *string    `json:"destination_account"`
	Amount              *int64     `json:"amount"`
	BalanceBefore       *int64     `json:"balance_before"`
	BalanceAfter        *int64     `json:"balance_after"`
	Type                *string    `json:"type"`
	Status              *string    `json:"status"`
}

type TransactionReponse struct {
	UserId              uuid.UUID `json:"user_id"`
	BankId              uuid.UUID `json:"bank_id"`
	TransactionId       string    `json:"transaction_id"`
	Destination_Account string    `json:"destination_account"`
	Amount              int64     `json:"amount"`
	BalanceBefore       int64     `json:"balance_before"`
	BalanceAfter        int64     `json:"balance_after"`
	Type                string    `json:"type"`
	Status              string    `json:"status"`
}

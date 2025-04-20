package dto

import "github.com/google/uuid"

type UserRequest struct {
	Name     *string         `json:"name"`
	Email    *string         `json:"email" binding:"omitempty,email"`
	Password *string         `json:"password"`
	Type     *string         `json:"type"`
	Phone    *string         `json:"phone"`
	Accounts *AccountRequest `json:"accounts"`
}

type AccountRequest struct {
	BankId        *uuid.UUID `json:"bank_id"`
	AccountNumber *string    `json:"account_number"`
	Balance       *int64     `json:"balance"`
	Status        *string    `json:"status"`
}

type UserResponse struct {
	ID       uuid.UUID        `json:"id"`
	Name     string           `json:"name"`
	Email    string           `json:"email" binding:"omitempty,email"`
	Password string           `json:"password"`
	Type     string           `json:"type"`
	Phone    string           `json:"phone"`
	Accounts *AccountResponse `json:"accounts"`
}

type AccountResponse struct {
	BankId        uuid.UUID `json:"bank_id"`
	AccountNumber string    `json:"account_number"`
	Balance       int64     `json:"balance"`
	Status        string    `json:"status"`
}

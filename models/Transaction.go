package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                  uuid.UUID `gorm:"type:UUID;default:uuid_generate_v4();primary_key"`
	UserId              uuid.UUID `gorm:"type:uuid;not null;"`
	BankId              uuid.UUID `gorm:"type:uuid;not null;" `
	TransactionId       string    `gorm:"unique" `
	Destination_Account string
	Amount              int64
	BalanceBefore       int64
	BalanceAfter        int64
	Type                string
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *gorm.DeletedAt `gorm:"index"`
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId        uuid.UUID `gorm:"type:uuid;not null;"`
	BankId        uuid.UUID `gorm:"type:uuid;not null"`
	AccountNumber string    `gorm:"unique"`
	Balance       int64
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *gorm.DeletedAt `gorm:"index"`

	// Relations
	User *User `gorm:"foreignKey:UserId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Bank Bank  `gorm:"foreignKey:BankId"`
}

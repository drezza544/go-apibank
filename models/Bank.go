package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bank struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Code            string    `gorm:"unique"`
	Name            string
	Address         string
	CostTransaction int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *gorm.DeletedAt `gorm:"index"`

	// Relations
	Accounts []Account `gorm:"foreignKey:BankId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

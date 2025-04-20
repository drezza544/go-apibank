package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	Type      string
	Phone     string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`
	Accounts  Account         `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

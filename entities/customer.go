package entities

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID               uuid.UUID `gorm:"primaryKey;autoIncrement:false"`
	Name             string    `gorm:"not null"`
	Email            string    `gorm:"uniqueIndex;size:320;not null"`
	Password         string    `gorm:"not null"`
	Photo            string    `gorm:"not null"`
	VerificationCode string    `gorm:"not null"`
	Verified         bool      `gorm:"default:false"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

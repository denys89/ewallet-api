package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	FirstName   string    `json:"first_name" gorm:"not null"`
	LastName    string    `json:"last_name" gorm:"not null"`
	PhoneNumber string    `json:"phone_number" gorm:"unique;not null"`
	Address     string    `json:"address" gorm:"not null"`
	Pin         string    `json:"-" gorm:"not null"`
	Balance     float64   `json:"balance" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

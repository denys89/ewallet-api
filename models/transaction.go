package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TOPUP, TRANSFER, PAYMENT, SUCCESS, DEBIT, CREDIT string = "TOPUP", "TRANSFER", "PAYMENT", "SUCCESS", "DEBIT", "CREDIT"
)

type Transaction struct {
	ID              uuid.UUID  `json:"id" gorm:"type:char(36);primary_key"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:char(36);not null"`
	Type            string     `json:"type" gorm:"not null"`
	TransactionType string     `json:"transaction_type" gorm:"not null"`
	Amount          float64    `json:"amount" gorm:"not null"`
	BalanceBefore   float64    `json:"balance_before" gorm:"not null"`
	BalanceAfter    float64    `json:"balance_after" gorm:"not null"`
	RecipientID     *uuid.UUID `json:"recipient_id,omitempty" gorm:"type:char(36)"`
	Description     string     `json:"description"`
	ReferenceNumber string     `json:"reference_number" gorm:"unique;not null"`
	Status          string     `json:"status" gorm:"not null"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	User            User       `json:"-" gorm:"foreignKey:UserID"`
	Recipient       *User      `json:"-" gorm:"foreignKey:RecipientID"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

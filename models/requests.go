package models

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	Address     string `json:"address" binding:"required"`
	Pin         string `json:"pin" binding:"required,numeric,len=6"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	Pin         string `json:"pin" binding:"required,numeric,len=6"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	Address     string `json:"address" binding:"required"`
}

type TransactionRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
}

type TransferRequest struct {
	RecipientID uuid.UUID `json:"recipient_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Description string    `json:"description" binding:"required"`
}

type TransactionFilter struct {
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	Type      string    `form:"type" binding:"omitempty,oneof=TOP_UP PAYMENT TRANSFER"`
}

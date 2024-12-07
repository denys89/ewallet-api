package repositories

import (
	"errors"

	"github.com/denys89/ewallet-api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPhoneNumberExists  = errors.New("phone number already registered")
	ErrInvalidCredentials = errors.New("phone number and pin doesn't match")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	// Check if phone number already exists
	exists := &models.User{}
	err := r.db.Where("phone_number = ?", user.PhoneNumber).First(exists).Error
	if err == nil {
		return ErrPhoneNumberExists
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return r.db.Create(user).Error
}

func (r *UserRepository) FindByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

package repositories

import (
	"github.com/denys89/ewallet-api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) getUserForUpdate(tx *gorm.DB, userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *TransactionRepository) GetUserTransactions(userID uuid.UUID, page, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	offset := (page - 1) * limit

	err := r.db.Where("user_id = ? OR recipient_id = ?", userID, userID).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) TopUp(userID uuid.UUID, amount float64) (uuid.UUID, float64, float64, error) {
	var transactionID uuid.UUID
	var balanceBefore, balanceAfter float64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		user, err := r.getUserForUpdate(tx, userID)
		if err != nil {
			return err
		}

		balanceBefore = user.Balance
		balanceAfter = balanceBefore + amount

		// Update user balance
		if err := tx.Model(user).Update("balance", balanceAfter).Error; err != nil {
			return err
		}

		// Create transaction record
		transaction := &models.Transaction{
			ID:              uuid.New(),
			UserID:          userID,
			Type:            models.CREDIT,
			TransactionType: models.TOPUP,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    balanceAfter,
			Amount:          amount,
			Status:          models.SUCCESS,
			Description:     "Top up balance",
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		transactionID = transaction.ID
		return nil
	})

	if err != nil {
		return uuid.UUID{}, 0, 0, err
	}

	return transactionID, balanceBefore, balanceAfter, nil
}

func (r *TransactionRepository) Payment(userID uuid.UUID, amount float64, remarks string) (*models.Transaction, float64, float64, error) {
	var transaction models.Transaction
	var balanceBefore, balanceAfter float64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		user, err := r.getUserForUpdate(tx, userID)
		if err != nil {
			return err
		}

		balanceBefore = user.Balance
		if balanceBefore < amount {
			return models.ErrInvalidTransaction
		}

		balanceAfter = balanceBefore - amount

		// Update user balance
		if err := tx.Model(user).Update("balance", balanceAfter).Error; err != nil {
			return err
		}

		// Create transaction record
		transaction = models.Transaction{
			ID:              uuid.New(),
			UserID:          userID,
			Type:            models.DEBIT,
			TransactionType: models.PAYMENT,
			Amount:          amount,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    balanceAfter,
			Description:     remarks,
			Status:          models.SUCCESS,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, 0, 0, err
	}

	return &transaction, balanceBefore, balanceAfter, nil
}

func (r *TransactionRepository) Transfer(userID uuid.UUID, amount float64, targetUser string, remarks string) (*models.Transaction, float64, float64, error) {
	var transaction models.Transaction
	var balanceBefore, balanceAfter float64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get sender with lock
		sender, err := r.getUserForUpdate(tx, userID)
		if err != nil {
			return err
		}

		// Get recipient with lock
		recipient, err := r.getUserForUpdate(tx, uuid.MustParse(targetUser))
		if err != nil {
			return err
		}

		balanceBefore = sender.Balance
		if balanceBefore < amount {
			return models.ErrInvalidTransaction
		}

		balanceAfter = balanceBefore - amount

		// Update sender's balance
		if err := tx.Model(sender).Update("balance", balanceAfter).Error; err != nil {
			return err
		}

		// Create sender's transaction record
		senderTransID := uuid.New()
		transaction = models.Transaction{
			ID:              senderTransID,
			UserID:          userID,
			Type:            models.DEBIT,
			TransactionType: models.TRANSFER,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    balanceAfter,
			Amount:          amount,
			Status:          models.SUCCESS,
			Description:     remarks,
			RecipientID:     &recipient.ID,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Update recipient's balance
		recipientBalanceBefore := recipient.Balance
		recipientBalanceAfter := recipientBalanceBefore + amount

		if err := tx.Model(recipient).Update("balance", recipientBalanceAfter).Error; err != nil {
			return err
		}

		// Create recipient's transaction record
		recipientTrans := models.Transaction{
			ID:              uuid.New(),
			UserID:          recipient.ID,
			Type:            models.CREDIT,
			TransactionType: models.TRANSFER,
			BalanceBefore:   recipientBalanceBefore,
			BalanceAfter:    recipientBalanceAfter,
			Amount:          amount,
			Status:          models.SUCCESS,
			ReferenceNumber: senderTransID.String(),
			Description:     remarks,
		}

		if err := tx.Create(&recipientTrans).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, 0, 0, err
	}

	return &transaction, balanceBefore, balanceAfter, nil
}

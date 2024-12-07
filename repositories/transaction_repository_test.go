package repositories

import (
	"testing"

	"github.com/denys89/ewallet-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository *TransactionRepository
	user       *models.User
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Transaction{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = &TransactionRepository{db: db}

	// Create a test user
	suite.user = &models.User{
		ID:          uuid.New(),
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		Pin:         "123456",
		Balance:     1000,
	}
	err = db.Create(suite.user).Error
	assert.NoError(suite.T(), err)
}

func (suite *TransactionRepositoryTestSuite) TestTopUp() {
	amount := float64(500)
	transactionID, balanceBefore, balanceAfter, err := suite.repository.TopUp(suite.user.ID, amount)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transactionID)
	assert.Equal(suite.T(), float64(1000), balanceBefore)
	assert.Equal(suite.T(), float64(1500), balanceAfter)

	// Verify transaction was created
	var transaction models.Transaction
	err = suite.db.First(&transaction, "id = ?", transactionID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "TOPUP", transaction.Type)
	assert.Equal(suite.T(), amount, transaction.Amount)

	// Verify user balance was updated
	var user models.User
	err = suite.db.First(&user, "id = ?", suite.user.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(1500), user.Balance)
}

func (suite *TransactionRepositoryTestSuite) TestPayment() {
	amount := float64(300)
	transactionID, balanceBefore, balanceAfter, err := suite.repository.Payment(suite.user.ID, amount, "Test payment")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transactionID)
	assert.Equal(suite.T(), float64(1000), balanceBefore)
	assert.Equal(suite.T(), float64(700), balanceAfter)

	// Verify transaction was created
	var transaction models.Transaction
	err = suite.db.First(&transaction, "id = ?", transactionID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "PAYMENT", transaction.Type)
	assert.Equal(suite.T(), amount, transaction.Amount)

	// Verify user balance was updated
	var user models.User
	err = suite.db.First(&user, "id = ?", suite.user.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(700), user.Balance)
}

func (suite *TransactionRepositoryTestSuite) TestGetUserTransactions() {
	// Create some test transactions
	transactions := []models.Transaction{
		{
			ID:            uuid.New(),
			UserID:        suite.user.ID,
			Type:          "TOPUP",
			Amount:        500,
			BalanceBefore: 1000,
			BalanceAfter:  1500,
		},
		{
			ID:            uuid.New(),
			UserID:        suite.user.ID,
			Type:          "PAYMENT",
			Amount:        200,
			BalanceBefore: 1500,
			BalanceAfter:  1300,
		},
	}

	for _, tx := range transactions {
		err := suite.db.Create(&tx).Error
		assert.NoError(suite.T(), err)
	}

	// Test pagination
	result, err := suite.repository.GetUserTransactions(suite.user.ID, 1, 10)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func (suite *TransactionRepositoryTestSuite) TestPaymentInsufficientBalance() {
	amount := float64(2000) // More than current balance
	_, _, _, err := suite.repository.Payment(suite.user.ID, amount, "Test payment")
	assert.Error(suite.T(), err)

	// Verify user balance remains unchanged
	var user models.User
	err = suite.db.First(&user, "id = ?", suite.user.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(1000), user.Balance)
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

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

type UserRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository *UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Migrate the schema
	err = db.AutoMigrate(&models.User{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = &UserRepository{db: db}
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	user := &models.User{
		ID:          uuid.New(),
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		Pin:         "123456",
		Balance:     0,
	}

	err := suite.repository.Create(user)
	assert.NoError(suite.T(), err)

	// Verify user was created
	var found models.User
	err = suite.db.First(&found, "id = ?", user.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.FirstName, found.FirstName)
	assert.Equal(suite.T(), user.PhoneNumber, found.PhoneNumber)
}

func (suite *UserRepositoryTestSuite) TestFindByID() {
	user := &models.User{
		ID:          uuid.New(),
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		Pin:         "123456",
		Balance:     0,
	}

	err := suite.repository.Create(user)
	assert.NoError(suite.T(), err)

	found, err := suite.repository.FindByID(user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, found.ID)
	assert.Equal(suite.T(), user.FirstName, found.FirstName)
}

func (suite *UserRepositoryTestSuite) TestFindByPhoneNumber() {
	user := &models.User{
		ID:          uuid.New(),
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		Pin:         "123456",
		Balance:     0,
	}

	err := suite.repository.Create(user)
	assert.NoError(suite.T(), err)

	found, err := suite.repository.FindByPhoneNumber(user.PhoneNumber)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, found.ID)
	assert.Equal(suite.T(), user.PhoneNumber, found.PhoneNumber)
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
	user := &models.User{
		ID:          uuid.New(),
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		Pin:         "123456",
		Balance:     0,
	}

	err := suite.repository.Create(user)
	assert.NoError(suite.T(), err)

	user.FirstName = "Jane"
	err = suite.repository.Update(user)
	assert.NoError(suite.T(), err)

	found, err := suite.repository.FindByID(user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Jane", found.FirstName)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

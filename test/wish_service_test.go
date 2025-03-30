package test

import (
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wishlist-app/internal/models"
	"wishlist-app/internal/service"
)

type MockWishRepository struct {
	mock.Mock
}

func (m *MockWishRepository) Create(wish *models.Wish) error {
	args := m.Called(wish)
	return args.Error(0)
}

func (m *MockWishRepository) GetByID(id uint) (*models.Wish, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Wish), args.Error(1)
}

func (m *MockWishRepository) Update(wish *models.Wish) error {
	args := m.Called(wish)
	return args.Error(0)
}

func (m *MockWishRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWishRepository) GetByUserID(userID uint) ([]models.Wish, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Wish), args.Error(1)
}

func (m *MockWishRepository) GetByUsername(username string) ([]models.Wish, error) {
	args := m.Called(username)
	return args.Get(0).([]models.Wish), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByLogin(login string) (*models.User, error) {
	args := m.Called(login)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Exists(login string) (bool, error) {
	args := m.Called(login)
	return args.Bool(0), args.Error(1)
}

func TestWishService_Create(t *testing.T) {
	wishService := service.NewWishService(mockWishRepo, mockUserRepo)

	testWish := &models.Wish{
		UserID: 1,
		Title:  "Test Wish",
	}

	mockWishRepo.On("Create", testWish).Return(nil)

	createdWish, err := wishService.Create(1, testWish)
	assert.NoError(t, err)
	assert.Equal(t, testWish, createdWish)
	mockWishRepo.AssertExpectations(t)
}

func TestWishService_GetByID(t *testing.T) {
	wishService := service.NewWishService(mockWishRepo, mockUserRepo)

	testWish := &models.Wish{
		Model:  gorm.Model{ID: 1, CreatedAt: time.Now()},
		UserID: 1,
		Title:  "Test Wish",
	}

	mockWishRepo.On("GetByID", uint(1)).Return(testWish, nil)

	wish, err := wishService.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, testWish, wish)
	mockWishRepo.AssertExpectations(t)
}

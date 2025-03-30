package repository

import (
	"time"

	"wishlist-app/internal/models"
	"wishlist-app/pkg/metrics"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByLogin(login string) (*models.User, error)
	Exists(login string) (bool, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	start := time.Now()
	err := r.db.Create(user).Error
	metrics.RecordDatabaseQuery("insert", "users", time.Since(start).Seconds())
	return err
}

func (r *UserRepository) FindByLogin(login string) (*models.User, error) {
	start := time.Now()
	var user models.User
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		metrics.RecordDatabaseQuery("select", "users", time.Since(start).Seconds())
		return nil, err
	}
	metrics.RecordDatabaseQuery("select", "users", time.Since(start).Seconds())
	return &user, nil
}

func (r *UserRepository) Exists(login string) (bool, error) {
	start := time.Now()
	var count int64
	err := r.db.Model(&models.User{}).Where("login = ?", login).Count(&count).Error
	metrics.RecordDatabaseQuery("select", "users", time.Since(start).Seconds())
	return count > 0, err
}

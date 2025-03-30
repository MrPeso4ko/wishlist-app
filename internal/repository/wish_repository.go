package repository

import (
	"wishlist-app/internal/models"

	"gorm.io/gorm"
)

type WishRepositoryInterface interface {
	Create(wish *models.Wish) error
	GetByID(id uint) (*models.Wish, error)
	Update(wish *models.Wish) error
	Delete(id uint) error
	GetByUserID(userID uint) ([]models.Wish, error)
	GetByUsername(username string) ([]models.Wish, error)
}

type WishRepository struct {
	db *gorm.DB
}

func NewWishRepository(db *gorm.DB) *WishRepository {
	return &WishRepository{db: db}
}

func (r *WishRepository) Create(wish *models.Wish) error {
	return r.db.Create(wish).Error
}

func (r *WishRepository) GetByID(id uint) (*models.Wish, error) {
	var wish models.Wish
	if err := r.db.Preload("User").First(&wish, id).Error; err != nil {
		return nil, err
	}
	return &wish, nil
}

func (r *WishRepository) Update(wish *models.Wish) error {
	return r.db.Save(wish).Error
}

func (r *WishRepository) Delete(id uint) error {
	return r.db.Delete(&models.Wish{}, id).Error
}

func (r *WishRepository) GetByUserID(userID uint) ([]models.Wish, error) {
	var wishes []models.Wish
	if err := r.db.Preload("User").Where("user_id = ?", userID).Find(&wishes).Error; err != nil {
		return nil, err
	}
	return wishes, nil
}

func (r *WishRepository) GetByUsername(username string) ([]models.Wish, error) {
	var wishes []models.Wish
	if err := r.db.Joins("User").Where("users.login = ?", username).Find(&wishes).Error; err != nil {
		return nil, err
	}
	return wishes, nil
}

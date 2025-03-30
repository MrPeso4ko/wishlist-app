package service

import (
	"errors"

	"wishlist-app/internal/models"
	"wishlist-app/internal/repository"
)

type WishService struct {
	wishRepo repository.WishRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

func NewWishService(wishRepo repository.WishRepositoryInterface, userRepo repository.UserRepositoryInterface) *WishService {
	return &WishService{
		wishRepo: wishRepo,
		userRepo: userRepo,
	}
}

func (s *WishService) Create(userID uint, wish *models.Wish) (*models.Wish, error) {
	wish.UserID = userID
	if err := s.wishRepo.Create(wish); err != nil {
		return nil, err
	}
	return wish, nil
}

func (s *WishService) GetByID(userID, wishID uint) (*models.Wish, error) {
	wish, err := s.wishRepo.GetByID(wishID)
	if err != nil {
		return nil, err
	}

	if wish.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return wish, nil
}

func (s *WishService) Update(userID uint, wish *models.Wish) error {
	existingWish, err := s.wishRepo.GetByID(wish.ID)
	if err != nil {
		return err
	}

	if existingWish.UserID != userID {
		return errors.New("unauthorized")
	}

	wish.UserID = userID
	return s.wishRepo.Update(wish)
}

func (s *WishService) Delete(userID, wishID uint) error {
	wish, err := s.wishRepo.GetByID(wishID)
	if err != nil {
		return err
	}

	if wish.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.wishRepo.Delete(wishID)
}

func (s *WishService) GetByUserID(userID uint) ([]models.Wish, error) {
	return s.wishRepo.GetByUserID(userID)
}

func (s *WishService) GetByUsername(username string) ([]models.Wish, error) {
	return s.wishRepo.GetByUsername(username)
}

package handler

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"wishlist-app/internal/config"
	"wishlist-app/internal/models"
	"wishlist-app/internal/service"
	"wishlist-app/pkg/logger"
	"wishlist-app/pkg/metrics"

	"github.com/gin-gonic/gin"
)

type WishHandler struct {
	wishService *service.WishService
	logger      logger.Logger
	cfg         *config.Config
}

func NewWishHandler(cfg *config.Config, logger logger.Logger, wishService *service.WishService) *WishHandler {
	return &WishHandler{
		wishService: wishService,
		cfg:         cfg,
		logger:      logger,
	}
}

type CreateWishRequest struct {
	Title    string  `json:"title" binding:"required"`
	Comment  string  `json:"comment"`
	ImageURL string  `json:"image_url"`
	Price    float64 `json:"price"`
}

type UpdateWishRequest struct {
	Title    string  `json:"title"`
	Comment  string  `json:"comment"`
	ImageURL string  `json:"image_url"`
	Price    float64 `json:"price"`
}

// Create godoc
// @Summary Create a new wish
// @Description Create a new wish for the authenticated user
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateWishRequest true "Create Wish Request"
// @Success 201 {object} models.PublicWish "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /wishes [post]
func (h *WishHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.RecordWishOperation("create", "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wish := &models.Wish{
		UserID:   userID,
		Title:    req.Title,
		Comment:  req.Comment,
		ImageURL: req.ImageURL,
		Price:    req.Price,
	}

	createdWish, err := h.wishService.Create(userID, wish)
	if err != nil {
		metrics.RecordWishOperation("create", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.RecordWishOperation("create", "success")
	c.JSON(http.StatusCreated, createdWish.ToPublic())
}

// Update godoc
// @Summary Update a wish
// @Description Update an existing wish for the authenticated user
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Wish ID"
// @Param request body UpdateWishRequest true "Update Wish Request"
// @Success 200 "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /wishes/{id} [put]
func (h *WishHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	wishID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		metrics.RecordWishOperation("update", "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wish ID"})
		return
	}

	var req UpdateWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.RecordWishOperation("update", "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wish := &models.Wish{
		Model:    gorm.Model{ID: uint(wishID)},
		Title:    req.Title,
		Comment:  req.Comment,
		ImageURL: req.ImageURL,
		Price:    req.Price,
	}

	if err := h.wishService.Update(userID, wish); err != nil {
		metrics.RecordWishOperation("update", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.RecordWishOperation("update", "success")
	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary Delete a wish
// @Description Delete an existing wish for the authenticated user
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Wish ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /wishes/{id} [delete]
func (h *WishHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	wishID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		metrics.RecordWishOperation("delete", "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wish ID"})
		return
	}

	if err := h.wishService.Delete(userID, uint(wishID)); err != nil {
		metrics.RecordWishOperation("delete", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.RecordWishOperation("delete", "success")
	c.Status(http.StatusNoContent)
}

// GetByUserID godoc
// @Summary Get wishes for authenticated user
// @Description Get all wishes for the authenticated user
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.PublicWish "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /wishes [get]
func (h *WishHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")

	wishes, err := h.wishService.GetByUserID(userID)
	if err != nil {
		metrics.RecordWishOperation("read", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.RecordWishOperation("read", "success")
	publicWishes := make([]*models.PublicWish, len(wishes))
	for i, wish := range wishes {
		publicWishes[i] = wish.ToPublic()
	}

	c.JSON(http.StatusOK, publicWishes)
}

// GetByUsername godoc
// @Summary Get wishes by username
// @Description Get all public wishes for a specific user by username
// @Tags wishes
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {array} models.PublicWish "OK"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /wishes/{username} [get]
func (h *WishHandler) GetByUsername(c *gin.Context) {
	username := c.Param("username")

	wishes, err := h.wishService.GetByUsername(username)
	if err != nil {
		metrics.RecordWishOperation("read", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.RecordWishOperation("read", "success")
	publicWishes := make([]*models.PublicWish, len(wishes))
	for i, wish := range wishes {
		publicWishes[i] = wish.ToPublic()
	}

	c.JSON(http.StatusOK, publicWishes)
}

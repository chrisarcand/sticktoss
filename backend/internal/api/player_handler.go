package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sticktoss/backend/internal/auth"
	"github.com/sticktoss/backend/internal/models"
	"gorm.io/gorm"
)

type PlayerHandler struct {
	db *gorm.DB
}

func NewPlayerHandler(db *gorm.DB) *PlayerHandler {
	return &PlayerHandler{db: db}
}

type CreatePlayerRequest struct {
	Name        string `json:"name" binding:"required"`
	SkillWeight int    `json:"skill_weight" binding:"required,min=1,max=5"`
}

type UpdatePlayerRequest struct {
	Name        string `json:"name"`
	SkillWeight int    `json:"skill_weight" binding:"omitempty,min=1,max=5"`
}

// GetPlayers returns all players for the authenticated user
func (h *PlayerHandler) GetPlayers(c *gin.Context) {
	userID := auth.GetUserID(c)

	var players []models.Player
	if err := h.db.Where("user_id = ?", userID).Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetPlayer returns a single player by ID
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	userID := auth.GetUserID(c)
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	var player models.Player
	if err := h.db.Where("id = ? AND user_id = ?", playerID, userID).First(&player).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

// CreatePlayer creates a new player
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	userID := auth.GetUserID(c)

	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player := models.Player{
		UserID:      userID,
		Name:        req.Name,
		SkillWeight: req.SkillWeight,
	}

	if err := h.db.Create(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create player"})
		return
	}

	c.JSON(http.StatusCreated, player)
}

// UpdatePlayer updates an existing player
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	userID := auth.GetUserID(c)
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	var player models.Player
	if err := h.db.Where("id = ? AND user_id = ?", playerID, userID).First(&player).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	var req UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		player.Name = req.Name
	}
	if req.SkillWeight > 0 {
		player.SkillWeight = req.SkillWeight
	}

	if err := h.db.Save(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update player"})
		return
	}

	c.JSON(http.StatusOK, player)
}

// DeletePlayer deletes a player
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	userID := auth.GetUserID(c)
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	var player models.Player
	if err := h.db.Where("id = ? AND user_id = ?", playerID, userID).First(&player).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	// Delete player (this will also remove from groups due to foreign key constraints)
	if err := h.db.Delete(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "player deleted"})
}

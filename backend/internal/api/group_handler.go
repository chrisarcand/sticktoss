package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sticktoss/backend/internal/auth"
	"github.com/sticktoss/backend/internal/models"
	"github.com/sticktoss/backend/internal/teamgen"
	"gorm.io/gorm"
)

type GroupHandler struct {
	db *gorm.DB
}

func NewGroupHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddPlayerToGroupRequest struct {
	PlayerID uint `json:"player_id" binding:"required"`
}

type GenerateTeamsRequest struct {
	NumTeams      int    `json:"num_teams" binding:"required,min=2"`
	LockedPlayers [][]uint `json:"locked_players"` // Array of arrays, each inner array is players that should be on same team
}

// GetGroups returns all groups for the authenticated user
func (h *GroupHandler) GetGroups(c *gin.Context) {
	userID := auth.GetUserID(c)

	var groups []models.Group
	if err := h.db.Where("user_id = ?", userID).Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroup returns a single group by ID with its players
func (h *GroupHandler) GetGroup(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var group models.Group
	if err := h.db.Preload("Players").Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	userID := auth.GetUserID(c)

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := models.Group{
		UserID: userID,
		Name:   req.Name,
	}

	if err := h.db.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateGroup updates an existing group
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var group models.Group
	if err := h.db.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group.Name = req.Name

	if err := h.db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteGroup deletes a group
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var group models.Group
	if err := h.db.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	if err := h.db.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}

// AddPlayerToGroup adds a player to a group
func (h *GroupHandler) AddPlayerToGroup(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var req AddPlayerToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify group belongs to user
	var group models.Group
	if err := h.db.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	// Verify player belongs to user
	var player models.Player
	if err := h.db.Where("id = ? AND user_id = ?", req.PlayerID, userID).First(&player).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	// Add player to group (GORM handles the many-to-many)
	if err := h.db.Model(&group).Association("Players").Append(&player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add player to group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "player added to group"})
}

// RemovePlayerFromGroup removes a player from a group
func (h *GroupHandler) RemovePlayerFromGroup(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	playerID, err := strconv.ParseUint(c.Param("player_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	// Verify group belongs to user
	var group models.Group
	if err := h.db.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	// Verify player belongs to user
	var player models.Player
	if err := h.db.Where("id = ? AND user_id = ?", playerID, userID).First(&player).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	// Remove player from group
	if err := h.db.Model(&group).Association("Players").Delete(&player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove player from group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "player removed from group"})
}

// GenerateTeams generates balanced teams for a group
func (h *GroupHandler) GenerateTeams(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var req GenerateTeamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get group with players
	var group models.Group
	if err := h.db.Preload("Players").Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	if len(group.Players) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group has no players"})
		return
	}

	if len(group.Players) < req.NumTeams {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough players for the requested number of teams"})
		return
	}

	// Generate teams
	teams, err := teamgen.GenerateBalancedTeams(group.Players, req.NumTeams, req.LockedPlayers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
	})
}

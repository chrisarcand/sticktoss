package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sticktoss/backend/internal/auth"
	"github.com/sticktoss/backend/internal/models"
	"github.com/sticktoss/backend/internal/teamgen"
	"github.com/sticktoss/backend/internal/utils"
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
	NumTeams         int      `json:"num_teams" binding:"required,min=2"`
	LockedPlayers    [][]uint `json:"locked_players"`     // Array of arrays, each inner array is players that should be on same team
	SeparatedPlayers [][]uint `json:"separated_players"`  // Array of arrays, each inner array is players that should be on different teams
	UseJerseyColors  bool     `json:"use_jersey_colors"`  // Whether to use jersey colors (Light/Dark)
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
	teams, err := teamgen.GenerateBalancedTeams(group.Players, req.NumTeams, req.LockedPlayers, req.SeparatedPlayers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate share ID for the game
	shareID, err := utils.GenerateShareID(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate share ID"})
		return
	}

	// Marshal teams data to JSON
	teamsJSON, err := json.Marshal(teams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save game"})
		return
	}

	// Save game to database
	game := models.Game{
		ShareID:         shareID,
		UserID:          userID,
		GroupID:         uint(groupID),
		GroupName:       group.Name,
		GroupLogo:       group.Logo,       // Copy logo for public access
		LogoContentType: group.LogoContentType,
		NumTeams:        req.NumTeams,
		UseJerseyColors: req.UseJerseyColors,
		TeamsData:       teamsJSON,
		CreatedAt:       time.Now(),
	}

	if err := h.db.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teams":    teams,
		"share_id": shareID,
	})
}

// GetGame retrieves a game by share ID (public endpoint, no auth required)
func (h *GroupHandler) GetGame(c *gin.Context) {
	shareID := c.Param("shareId")

	var game models.Game
	if err := h.db.Where("share_id = ?", shareID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	// Unmarshal teams data
	var teams []teamgen.Team
	if err := json.Unmarshal(game.TeamsData, &teams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load game data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"share_id":          game.ShareID,
		"group_name":        game.GroupName,
		"num_teams":         game.NumTeams,
		"use_jersey_colors": game.UseJerseyColors,
		"teams":             teams,
		"created_at":        game.CreatedAt,
		"has_logo":          len(game.GroupLogo) > 0,
	})
}

// UploadGroupLogo handles logo upload for a group
func (h *GroupHandler) UploadGroupLogo(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	// Verify group belongs to user
	var group models.Group
	if err := h.db.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Validate file size (max 2MB)
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 2MB)"})
		return
	}

	// Validate content type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/svg+xml" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type (only PNG, JPG, SVG allowed)"})
		return
	}

	// Read file data
	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer fileData.Close()

	logoData := make([]byte, file.Size)
	if _, err := fileData.Read(logoData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file data"})
		return
	}

	// Update group with logo
	group.Logo = logoData
	group.LogoContentType = contentType

	if err := h.db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save logo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logo uploaded successfully"})
}

// GetGroupLogo retrieves a group's logo (public endpoint)
func (h *GroupHandler) GetGroupLogo(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var group models.Group
	if err := h.db.Select("logo", "logo_content_type").Where("id = ?", groupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	if len(group.Logo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no logo found"})
		return
	}

	c.Data(http.StatusOK, group.LogoContentType, group.Logo)
}

// GetGameLogo retrieves a game's logo (public endpoint)
func (h *GroupHandler) GetGameLogo(c *gin.Context) {
	shareID := c.Param("shareId")

	var game models.Game
	if err := h.db.Select("group_logo", "logo_content_type").Where("share_id = ?", shareID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if len(game.GroupLogo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no logo found"})
		return
	}

	c.Data(http.StatusOK, game.LogoContentType, game.GroupLogo)
}

// DeleteGroupLogo deletes a group's logo
func (h *GroupHandler) DeleteGroupLogo(c *gin.Context) {
	userID := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	// Verify group belongs to user
	var group models.Group
	if err := h.db.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	// Clear logo
	group.Logo = nil
	group.LogoContentType = ""

	if err := h.db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete logo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logo deleted successfully"})
}

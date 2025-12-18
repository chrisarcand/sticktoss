package models

import (
	"time"
)

type Game struct {
	ShareID         string    `gorm:"primaryKey;size:12" json:"share_id"`
	UserID          uint      `json:"user_id"`
	GroupID         uint      `json:"group_id"`
	GroupName       string    `gorm:"size:255" json:"group_name"`
	GroupLogo       []byte    `gorm:"type:bytea" json:"-"` // Denormalized logo for public access
	LogoContentType string    `gorm:"size:50" json:"logo_content_type,omitempty"`
	NumTeams        int       `json:"num_teams"`
	UseJerseyColors bool      `json:"use_jersey_colors"`
	TeamsData       []byte    `gorm:"type:jsonb" json:"teams_data"` // Stores the complete team assignments
	CreatedAt       time.Time `json:"created_at"`
}

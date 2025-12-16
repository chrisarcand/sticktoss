package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user account
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Players []Player `gorm:"foreignKey:UserID" json:"-"`
	Groups  []Group  `gorm:"foreignKey:UserID" json:"-"`
}

// Player represents a hockey player with a skill weight
type Player struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	SkillWeight int       `gorm:"not null;check:skill_weight >= 1 AND skill_weight <= 5" json:"skill_weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User   User    `gorm:"foreignKey:UserID" json:"-"`
	Groups []Group `gorm:"many2many:group_players;" json:"-"`
}

// Group represents a collection of players
type Group struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    User     `gorm:"foreignKey:UserID" json:"-"`
	Players []Player `gorm:"many2many:group_players;" json:"players,omitempty"`
}

// GroupPlayer is the junction table for the many-to-many relationship
type GroupPlayer struct {
	GroupID  uint `gorm:"primaryKey" json:"group_id"`
	PlayerID uint `gorm:"primaryKey" json:"player_id"`
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Player{}, &Group{}, &GroupPlayer{})
}

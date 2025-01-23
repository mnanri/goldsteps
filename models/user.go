package models

import "time"

// User Model
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"` // Email address is unique
	Password string `json:"-"`                   // Password is private
	// Events    []Event   `json:"events" gorm:"foreignKey:UserID"` // Relationship between
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

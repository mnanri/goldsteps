package models

import (
	"errors"
	"math/rand"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User Model
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Alias     string    `json:"alias"`
	Password  string    `json:"password"` // Password is private
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate password (more than 8 characters by digits and alphabet)
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one number")
	}

	return nil
}

// BeforeCreate Hook: Generate a random 9-digit ID and hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a 9-digit random ID
	u.ID = uint(rand.Intn(900000000) + 100000000) // 100,000,000 ～ 999,999,999

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

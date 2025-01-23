package models

import (
	"errors"
	"time"
)

type Status string

const (
	ToDo       Status = "To Do"
	InProgress Status = "In Progress"
	Pending    Status = "Pending"
	InReview   Status = "In Review"
	Done       Status = "Done"
)

type Tag string

const (
	Urgent Tag = "Urgent"
	Medium Tag = "Medium"
	Low    Tag = "Low"
)

// Event Model
type Event struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	UserID      uint
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Deadline    time.Time `json:"deadline"`
	Status      Status    `json:"status" gorm:"type:text;default:'To Do'"`
	Tag         Tag       `json:"tag" gorm:"type:text;default:'Medium'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validation
func (e *Event) Validate() error {
	if !isValidStatus(e.Status) {
		return errors.New("invalid status value")
	}
	if !isValidTag(e.Tag) {
		return errors.New("invalid tag value")
	}
	return nil
}

func isValidStatus(s Status) bool {
	switch s {
	case ToDo, InProgress, Pending, InReview, Done:
		return true
	}
	return false
}

func isValidTag(t Tag) bool {
	switch t {
	case Urgent, Medium, Low:
		return true
	}
	return false
}

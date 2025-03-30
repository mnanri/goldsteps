package models

type NewsArticle struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"  gorm:"not null"`
	Link        string `json:"link" gorm:"unique;not null"`
	Description string `json:"description"`
}

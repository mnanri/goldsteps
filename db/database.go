package db

import (
	"goldsteps/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init DB
func InitDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open("steps.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migration of the tables
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.NewsArticle{},
		&models.Stock{},
		&models.StockDetail{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}

	return DB
}

package repository

import (
	"log"
	"server/models"

	"gorm.io/gorm"
)

func SaveNewsArticles(db *gorm.DB, articles []models.NewsArticle) error {
	for _, article := range articles {
		// Execude previous saved news
		var existing models.NewsArticle

		// Check existance
		err := db.Where("link = ?", article.Link).First(&existing).Error

		if err != nil {
			// Add new data
			if err == gorm.ErrRecordNotFound {
				if saveErr := db.Create(&article).Error; saveErr != nil {
					log.Println("Failed to save news:", saveErr)
					return saveErr
				}
			} else {
				// Unpredicted error
				log.Println("Error checking existing news:", err)
				return err
			}
		}
	}
	return nil
}

// Get all saved news
func GetAllNewsArticles(db *gorm.DB) ([]models.NewsArticle, error) {
	var articles []models.NewsArticle
	err := db.Find(&articles).Error
	return articles, err
}

// Search articles
func SearchNewsArticles(db *gorm.DB, query string) ([]models.NewsArticle, error) {
	var articles []models.NewsArticle

	// DEBUG
	// fmt.Println("Query: ", query)

	// Query
	err := db.Where("title LIKE ?", "%"+query+"%").
		Find(&articles).Error

	return articles, err
}

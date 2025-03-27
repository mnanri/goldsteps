package repository

import (
	"goldsteps/db"
	"goldsteps/models"
)

func GetAllMilestones() ([]models.Milestone, error) {
	var items []models.Milestone
	if err := db.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func CreateMilestone(item *models.Milestone) error {
	if err := db.DB.Create(item).Error; err != nil {
		return err
	}
	return nil
}

func DeleteMilestoneByID(id string) error {
	if err := db.DB.Delete(&models.Milestone{}, id).Error; err != nil {
		return err
	}
	return nil
}

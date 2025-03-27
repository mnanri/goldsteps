package repository

import (
	"goldsteps/db"
	"goldsteps/models"
)

func GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	if err := db.DB.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func GetEventByID(id string) (*models.Event, error) {
	var event models.Event
	if err := db.DB.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func CreateEvent(event *models.Event) error {
	return db.DB.Create(event).Error
}

func UpdateEvent(event *models.Event) error {
	return db.DB.Save(event).Error
}

func DeleteEvent(id string) error {
	return db.DB.Delete(&models.Event{}, id).Error
}

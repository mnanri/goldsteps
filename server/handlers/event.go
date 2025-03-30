package handlers

import (
	"net/http"
	"server/models"
	"server/repository"

	"github.com/labstack/echo/v4"
)

func GetEvents(c echo.Context) error {
	events, err := repository.GetAllEvents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch events"})
	}
	return c.JSON(http.StatusOK, events)
}

func GetEvent(c echo.Context) error {
	id := c.Param("id")
	event, err := repository.GetEventByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
	}
	return c.JSON(http.StatusOK, event)
}

func CreateEvent(c echo.Context) error {
	event := new(models.Event)
	if err := c.Bind(event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := repository.CreateEvent(event); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create event"})
	}

	return c.JSON(http.StatusCreated, event)
}

func UpdateEvent(c echo.Context) error {
	id := c.Param("id")
	event, err := repository.GetEventByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
	}

	updatedEvent := new(models.Event)
	if err := c.Bind(updatedEvent); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Validate status and tag
	if err := updatedEvent.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	event.Title = updatedEvent.Title
	event.Description = updatedEvent.Description
	event.StartTime = updatedEvent.StartTime
	event.EndTime = updatedEvent.EndTime
	event.Deadline = updatedEvent.Deadline
	event.Status = updatedEvent.Status
	event.Tag = updatedEvent.Tag

	if err := repository.UpdateEvent(event); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update event"})
	}

	return c.JSON(http.StatusOK, event)
}

func DeleteEvent(c echo.Context) error {
	id := c.Param("id")
	if err := repository.DeleteEvent(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete event"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Event deleted successfully"})
}

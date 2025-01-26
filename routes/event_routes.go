package routes

import (
	"fmt"
	"goldsteps/db"
	"goldsteps/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Event API
func RegisterEventRoutes(e *echo.Group) {
	// Get all events
	e.GET("/events", func(c echo.Context) error {
		// Debug: Print the request body
		// fmt.Println("Request Body:", c.Request().Body)

		var events []models.Event
		result := db.DB.Find(&events)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusOK, events)
	})

	// Get an event
	e.GET("/events/:id", func(c echo.Context) error {
		// Debug: Print the request body
		// fmt.Println("Request Body:", c.Request().Body)
		fmt.Println("Request Param:", c.Param("id"))

		id := c.Param("id")
		var event models.Event
		result := db.DB.First(&event, id)
		if result.Error != nil {
			if result.RowsAffected == 0 {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
			}
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusOK, event)
	})

	// Make an event
	e.POST("/events", func(c echo.Context) error {
		event := new(models.Event)
		if err := c.Bind(event); err != nil {
			fmt.Println("Bind Error:", err)
			return c.JSON(http.StatusBadRequest, err)
		}
		result := db.DB.Create(event)
		if result.Error != nil {
			fmt.Println("Database Error:", result.Error)
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusCreated, event)
	})

	// Update an event
	e.PUT("/events/:id", func(c echo.Context) error {
		id := c.Param("id")
		var event models.Event

		// Get the target event
		if err := db.DB.First(&event, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
		}

		// Update the data from the information from the request body
		updatedEvent := new(models.Event)
		if err := c.Bind(updatedEvent); err != nil {
			return c.JSON(http.StatusBadRequest, err)
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

		if err := db.DB.Save(&event).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, event)
	})

	// Delete an event
	e.DELETE("/events/:id", func(c echo.Context) error {
		id := c.Param("id")
		var event models.Event

		// Get the target event
		if err := db.DB.First(&event, id).Error; err != nil {
			if err.Error() == "record not found" {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
			}
			return c.JSON(http.StatusInternalServerError, err)
		}

		// Delete
		if err := db.DB.Delete(&event).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Event deleted successfully"})
	})
}

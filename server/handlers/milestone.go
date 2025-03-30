package handlers

import (
	"net/http"
	"server/models"
	"server/repository"

	"github.com/labstack/echo/v4"
)

func GetMilestones(c echo.Context) error {
	items, err := repository.GetAllMilestones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch milestone list"})
	}
	return c.JSON(http.StatusOK, items)
}

func AddMilestone(c echo.Context) error {
	item := new(models.Milestone)
	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := repository.CreateMilestone(item); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save item"})
	}

	return c.JSON(http.StatusCreated, item)
}

func DeleteMilestone(c echo.Context) error {
	id := c.Param("id")
	if err := repository.DeleteMilestoneByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete item"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Deleted successfully"})
}

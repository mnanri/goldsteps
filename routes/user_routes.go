package routes

import (
	"goldsteps/db"
	"goldsteps/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Group) {
	// Get a user
	e.GET("/user/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user models.User
		result := db.DB.First(&user, id)
		if result.Error != nil {
			if result.RowsAffected == 0 {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusOK, user)
	})

	// Make a user
	e.POST("/user", func(c echo.Context) error {
		user := new(models.User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		result := db.DB.Create(user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusCreated, user)
	})

	// Update a user
	e.PUT("/user/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user models.User

		// Get the target user
		if err := db.DB.First(&user, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		// Update the data from the information from the request body
		updatedUser := new(models.User)
		if err := c.Bind(updatedUser); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		user.Name = updatedUser.Name
		user.Email = updatedUser.Email

		if err := db.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, user)
	})

	// Delete a user
	e.DELETE("/user/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user models.User

		// Get the target user
		if err := db.DB.First(&user, id).Error; err != nil {
			if err.Error() == "record not found" {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, err)
		}

		// Delete
		if err := db.DB.Delete(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
	})
}

package routes

import (
	"fmt"
	"goldsteps/db"
	"goldsteps/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	// Read .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found, using environment variables or default value.")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("Warning: JWT_SECRET is not set, using default secret.")
		return "default_secret_key" // default for DEBUG
	}
	return secret
}

func RegisterUserRoutes(e *echo.Group) {
	// Get a user
	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user models.User
		result := db.DB.First(&user, id)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusOK, user)
	})

	// Register (Create) a user
	e.POST("/users", func(c echo.Context) error {
		user := new(models.User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if user.Alias == "" || user.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Alias and Password are required"})
		}
		if err := models.ValidatePassword(user.Password); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		result := db.DB.Create(user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
	})

	// Login
	e.POST("/login", func(c echo.Context) error {
		req := struct {
			Alias    string `json:"alias"`
			Password string `json:"password"`
		}{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		var user models.User
		if err := db.DB.Where("alias = ?", req.Alias).First(&user).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid alias or password"})
		}

		// Compare hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid alias or password"})
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
	})

	// Update a user (requires authentication)
	e.PUT("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user models.User

		if err := db.DB.First(&user, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		updatedUser := new(models.User)
		if err := c.Bind(updatedUser); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		if updatedUser.Alias != "" {
			user.Alias = updatedUser.Alias
		}

		if err := db.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		}

		return c.JSON(http.StatusOK, user)
	})

	// Delete a user
	e.DELETE("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user models.User

		if err := db.DB.First(&user, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		if err := db.DB.Delete(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
	})
}

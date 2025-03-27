package routes

import (
	"fmt"
	"goldsteps/db"
	"goldsteps/models"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found, using environment variables or default value.")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("Warning: JWT_SECRET is not set, using default secret.")
		return "default_secret_key"
	}
	return secret
}

func getUserIDFromToken(c echo.Context) (uint, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token")
	}
	return uint(userIDFloat), nil
}

func RegisterUserRoutes(e *echo.Group) {
	e.POST("/user/create", func(c echo.Context) error {
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}
		user.Password = string(hashedPassword)

		if err := db.DB.Create(user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
	})

	e.POST("/auth/login", func(c echo.Context) error {
		var req struct {
			Alias    string `json:"alias"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			// log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		var user models.User
		if err := db.DB.Where("alias = ?", req.Alias).First(&user).Error; err != nil {
			// log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid alias or password"})
		}

		log.Println("Stored Hash:", user.Password)
		log.Println("Input Password:", req.Password)

		if err := bcrypt.CompareHashAndPassword([]byte(strings.TrimSpace(user.Password)), []byte(strings.TrimSpace(req.Password))); err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid alias or password"})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
	})

	e.GET("/auth/ref", func(c echo.Context) error {
		userID, err := getUserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		return c.JSON(http.StatusOK, user)
	})

	e.POST("/auth/logout", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
	})

	e.PUT("/auth/conv", func(c echo.Context) error {
		userID, err := getUserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
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

	e.DELETE("/auth/rm", func(c echo.Context) error {
		userID, err := getUserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		if err := db.DB.Delete(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
	})

}

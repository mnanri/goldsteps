package main

import (
	"goldsteps/db"
	"goldsteps/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Init DB
	db.InitDB()

	// Root Endpoint
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome!")
	})

	e.GET("/api/data", func(c echo.Context) error {
		data := map[string]string{"message": "Hello from API"}
		return c.JSON(http.StatusOK, data)
	})

	// Set the routing
	api := e.Group("/api")
	routes.RegisterEventRoutes(api)
	routes.RegisterUserRoutes(api)

	// Awake server
	e.Logger.Fatal(e.Start(":8080"))
}

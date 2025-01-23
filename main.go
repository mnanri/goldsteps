package main

import (
	"goldsteps/db"
	"goldsteps/routes"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	e := echo.New()

	// Init DB
	db.InitDB()

	// Root Endpoint
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome!")
	})

	// Set the routing
	api := e.Group("/api")
	routes.RegisterEventRoutes(api)
	routes.RegisterUserRoutes(api)

	// Awake server
	e.Logger.Fatal(e.Start(":8080"))
}

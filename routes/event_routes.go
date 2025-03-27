package routes

import (
	"goldsteps/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterEventRoutes(e *echo.Group) {
	e.GET("/events", handlers.GetEvents)
	e.GET("/events/:id", handlers.GetEvent)
	e.POST("/events", handlers.CreateEvent)
	e.PUT("/events/:id", handlers.UpdateEvent)
	e.DELETE("/events/:id", handlers.DeleteEvent)
}

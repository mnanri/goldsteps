package routes

import (
	"goldsteps/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterMilestoneRoutes(e *echo.Group) {
	e.GET("/milestones", handlers.GetMilestones)
	e.POST("/milestones", handlers.AddMilestone)
	e.DELETE("/milestones/:id", handlers.DeleteMilestone)
}

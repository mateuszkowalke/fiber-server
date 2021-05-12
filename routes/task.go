package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mateuszkowalke/nozbe-tasks/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", controllers.Main)
	app.Post("login", controllers.Login)
	app.Get("/api/v1/tasks", controllers.GetTasks)
	app.Get("/api/v1/task/:id", controllers.GetTask)
	app.Post("/api/v1/task", controllers.NewTask)
	app.Put("/api/v1/task/:id", controllers.UpdateTask)
	app.Delete("/api/v1/task/:id", controllers.DeleteTask)
}

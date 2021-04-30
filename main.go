package main

import (
	"log"

	"github.com/mateuszkowalke/nozbe-tasks/database"
	"github.com/mateuszkowalke/nozbe-tasks/task"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/tasks", task.GetTasks)
	app.Get("/api/v1/task/:id", task.GetTask)
	app.Post("/api/v1/task", task.NewTask)
	app.Put("/api/v1/task/:id", task.UpdateTask)
	app.Delete("/api/v1/task/:id", task.DeleteTask)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))

}

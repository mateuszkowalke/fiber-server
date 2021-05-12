package main

import (
	"log"

	"github.com/mateuszkowalke/nozbe-tasks/database"
	"github.com/mateuszkowalke/nozbe-tasks/middleware"
	"github.com/mateuszkowalke/nozbe-tasks/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	app := fiber.New()
	app.Use(logger.New())

	app.Static("static", "front/dist")

	app.Use("/api", middleware.GetFromNozbe)

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))

}

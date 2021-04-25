package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mateuszkowalke/fiber-server/database"
	"github.com/mateuszkowalke/fiber-server/task"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/task", task.GetTasks)
	app.Get("/api/v1/task/:id", task.GetTask)
	app.Post("/api/v1/task", task.NewTask)
	app.Put("/api/v1/task/:id", task.UpdateTask)
	app.Delete("/api/v1/task/:id", task.DeleteTask)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "tasks.db")
	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("database connection opened")

	database.DBConn.AutoMigrate(&task.Task{})
	fmt.Println("database migrated")
}

func main() {

	app := fiber.New()

	initDatabase()
	defer database.DBConn.Close()

	setupRoutes(app)

	app.Listen(":3000")

}

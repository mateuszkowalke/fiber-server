package task

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/mateuszkowalke/fiber-server/database"
)

type Task struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

func GetTasks(c *fiber.Ctx) error {
	db := database.DBConn
	var tasks []Task
	db.Find(&tasks)
	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var task Task
	db.Find(&task, id)
	return c.JSON(task)
}

func NewTask(c *fiber.Ctx) error {
	db := database.DBConn
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&task)
	return c.JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
	return c.SendString("Update Task")
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var task Task
	db.First(&task, id)
	if task.Name == "" {
		c.Status(500).SendString("no task with given id")
	}
	db.Delete(&task)
	return c.SendString("successfully deleted task")
}

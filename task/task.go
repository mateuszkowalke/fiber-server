package task

import (
	"os"
	"time"

	"github.com/mateuszkowalke/nozbe-tasks/database"
	"github.com/mateuszkowalke/nozbe-tasks/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTasks(c *fiber.Ctx) error {
	taskCollection := database.MI.DB.Collection(os.Getenv("TASKS_COLLECTION"))

	query := bson.D{{}}

	cursor, err := taskCollection.Find(c.Context(), query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	var tasks []models.Task = make([]models.Task, 0)

	err = cursor.All(c.Context(), &tasks)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"tasks": tasks,
		},
	})
}

func GetTask(c *fiber.Ctx) error {
	taskCollection := database.MI.DB.Collection(os.Getenv("TASKS_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse Id",
			"error":   err,
		})
	}

	task := &models.Task{}

	query := bson.D{{Key: "_id", Value: id}}

	err = taskCollection.FindOne(c.Context(), query).Decode(task)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Task Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"task": task,
		},
	})
}

func NewTask(c *fiber.Ctx) error {
	taskCollection := database.MI.DB.Collection(os.Getenv("TASKS_COLLECTION"))

	data := new(models.Task)

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	data.ID = nil
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	result, err := taskCollection.InsertOne(c.Context(), data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot insert task",
			"error":   err,
		})
	}

	// get the inserted data
	task := &models.Task{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	taskCollection.FindOne(c.Context(), query).Decode(task)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"task": task,
		},
	})
}

func UpdateTask(c *fiber.Ctx) error {
	taskCollection := database.MI.DB.Collection(os.Getenv("TASKS_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	data := new(models.Task)
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	var dataToUpdate bson.D

	if data.Name != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "name", Value: data.Name})
	}

	if data.Description != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "description", Value: data.Description})
	}

	if data.Priority != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "priority", Value: data.Priority})
	}

	dataToUpdate = append(dataToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	update := bson.D{
		{Key: "$set", Value: dataToUpdate},
	}

	err = taskCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Task Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update task",
			"error":   err,
		})
	}

	task := &models.Task{}

	taskCollection.FindOne(c.Context(), query).Decode(task)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"task": task,
		},
	})
}

func DeleteTask(c *fiber.Ctx) error {
	taskCollection := database.MI.DB.Collection(os.Getenv("TASKS_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	err = taskCollection.FindOneAndDelete(c.Context(), query).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Task Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete task",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

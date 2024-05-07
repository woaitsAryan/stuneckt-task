package controllers

import (
	"github.com/woaitsAryan/stuneckt-task/initializers"
	"github.com/woaitsAryan/stuneckt-task/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("loggedinUser").(*models.User)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	initializers.DB.First(&user, "username = ?", username)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No user of this username found."}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"user":    user,
	})
}

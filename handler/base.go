package handler

import (
	"github.com/emamulandalib/airbringr-notification/response"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.JSON(response.Success{
		Message: "V1!",
	})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(response.Success{
		Message: "404!",
		Errors:  []string{"Path not found!"},
	})
}

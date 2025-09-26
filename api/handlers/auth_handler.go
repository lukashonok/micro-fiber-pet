package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MeHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// middleware кладёт user в c.Locals("user")
		val := c.Locals("user")
		if val == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		// val будет *entities.User
		return c.JSON(val)
	}
}

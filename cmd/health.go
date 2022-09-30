package main

import (
	"github.com/gofiber/fiber/v2"
)

// Health hanlder gets application health status on '/health' endpoint with
// status code 200 ok.
func (app *Application) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "available",
		"app":         app.Config.App,
		"version":     app.Config.Version,
		"environment": app.Config.Server.Env,
	})
}
